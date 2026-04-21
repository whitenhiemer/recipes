package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/whitenhiemer/recipe-site/internal/recipe"
	"github.com/whitenhiemer/recipe-site/internal/store"
)

type importURLRequest struct {
	URL      string `json:"url"`
	Category string `json:"category"`
}

func (s *Server) handleImportURL(w http.ResponseWriter, r *http.Request) {
	var req importURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" || req.Category == "" {
		http.Error(w, "url and category are required", http.StatusBadRequest)
		return
	}

	parsed, err := url.Parse(req.URL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}

	client := &http.Client{Timeout: 15 * time.Second}
	httpReq, err := http.NewRequest("GET", req.URL, nil)
	if err != nil {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	httpReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	httpReq.Header.Set("Accept-Language", "en-US,en;q=0.9")
	resp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, "failed to fetch URL: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		http.Error(w, fmt.Sprintf("URL returned status %d", resp.StatusCode), http.StatusBadGateway)
		return
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
	if err != nil {
		http.Error(w, "failed to read response", http.StatusBadGateway)
		return
	}

	html := string(body)
	rd := extractRecipeData(html)
	if rd.Title == "" {
		http.Error(w, "could not extract recipe from this URL — try pasting the recipe manually", http.StatusUnprocessableEntity)
		return
	}

	sourceDomain := parsed.Hostname()
	md := buildImportedMarkdown(rd, req.URL, sourceDomain)

	category := toSlug(req.Category)
	slug := toSlug(rd.Title)
	if slug == "" {
		http.Error(w, "could not generate slug from title", http.StatusUnprocessableEntity)
		return
	}

	parsedRecipe, err := recipe.ParseMarkdown([]byte(md), slug, category)
	if err != nil {
		http.Error(w, "failed to parse generated markdown", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	row := &store.RecipeRow{
		Slug:         slug,
		Title:        parsedRecipe.Title,
		Category:     category,
		Image:        parsedRecipe.Image,
		PrepTime:     parsedRecipe.PrepTime,
		CookTime:     parsedRecipe.CookTime,
		Servings:     parsedRecipe.Servings,
		Tags:         parsedRecipe.Tags,
		Ingredients:  toIngredientRows(parsedRecipe.Ingredients),
		Instructions: parsedRecipe.Instructions,
		Notes:        parsedRecipe.Notes,
		Markdown:     md,
		HTMLContent:  parsedRecipe.HTMLContent,
		SourceURL:    req.URL,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if _, err := s.db.InsertRecipe(row); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			http.Error(w, "a recipe with this name already exists in this category", http.StatusConflict)
			return
		}
		http.Error(w, "failed to save recipe", http.StatusInternalServerError)
		return
	}

	if err := s.rebuildIndex(); err != nil {
		http.Error(w, "recipe saved but index rebuild failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"slug":     slug,
		"category": category,
		"title":    rd.Title,
		"markdown": md,
	})
}

type recipeData struct {
	Title        string
	PrepTime     string
	CookTime     string
	Servings     string
	Ingredients  []string
	Instructions []string
}

var jsonLDPattern = regexp.MustCompile(`<script[^>]*type="application/ld\+json"[^>]*>([\s\S]*?)</script>`)

func extractRecipeData(html string) recipeData {
	if rd := extractFromJSONLD(html); rd.Title != "" {
		return rd
	}
	return recipeData{}
}

func extractFromJSONLD(html string) recipeData {
	matches := jsonLDPattern.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		jsonStr := match[1]

		// Try as single object
		if rd := parseJSONLDRecipe(jsonStr); rd.Title != "" {
			return rd
		}

		// Try as array
		var arr []json.RawMessage
		if json.Unmarshal([]byte(jsonStr), &arr) == nil {
			for _, item := range arr {
				if rd := parseJSONLDRecipe(string(item)); rd.Title != "" {
					return rd
				}
			}
		}

		// Try as @graph
		var graph struct {
			Graph []json.RawMessage `json:"@graph"`
		}
		if json.Unmarshal([]byte(jsonStr), &graph) == nil {
			for _, item := range graph.Graph {
				if rd := parseJSONLDRecipe(string(item)); rd.Title != "" {
					return rd
				}
			}
		}
	}
	return recipeData{}
}

func parseJSONLDRecipe(jsonStr string) recipeData {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		return recipeData{}
	}

	if !isRecipeType(obj["@type"]) {
		return recipeData{}
	}

	rd := recipeData{}
	rd.Title, _ = obj["name"].(string)
	rd.Title = strings.TrimSpace(rd.Title)

	if pt, ok := obj["prepTime"].(string); ok {
		rd.PrepTime = parseISODuration(pt)
	}
	if ct, ok := obj["cookTime"].(string); ok {
		rd.CookTime = parseISODuration(ct)
	}
	if yield, ok := obj["recipeYield"].(string); ok {
		rd.Servings = yield
	} else if yields, ok := obj["recipeYield"].([]interface{}); ok && len(yields) > 0 {
		rd.Servings = fmt.Sprintf("%v", yields[0])
	}

	if ings, ok := obj["recipeIngredient"].([]interface{}); ok {
		for _, ing := range ings {
			if s, ok := ing.(string); ok {
				rd.Ingredients = append(rd.Ingredients, strings.TrimSpace(s))
			}
		}
	}

	if insts, ok := obj["recipeInstructions"].([]interface{}); ok {
		for _, inst := range insts {
			switch v := inst.(type) {
			case string:
				rd.Instructions = append(rd.Instructions, strings.TrimSpace(v))
			case map[string]interface{}:
				if text, ok := v["text"].(string); ok {
					rd.Instructions = append(rd.Instructions, strings.TrimSpace(text))
				}
			}
		}
	}

	return rd
}

func isRecipeType(v interface{}) bool {
	if s, ok := v.(string); ok {
		return strings.EqualFold(s, "Recipe")
	}
	if arr, ok := v.([]interface{}); ok {
		for _, item := range arr {
			if s, ok := item.(string); ok && strings.EqualFold(s, "Recipe") {
				return true
			}
		}
	}
	return false
}

var isoDurationPattern = regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)

func parseISODuration(iso string) string {
	m := isoDurationPattern.FindStringSubmatch(iso)
	if m == nil {
		return iso
	}
	var parts []string
	if m[1] != "" {
		parts = append(parts, m[1]+" hr")
	}
	if m[2] != "" {
		parts = append(parts, m[2]+" min")
	}
	if m[3] != "" && len(parts) == 0 {
		parts = append(parts, m[3]+" sec")
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}

func buildImportedMarkdown(rd recipeData, sourceURL, sourceDomain string) string {
	var b strings.Builder

	b.WriteString("---\ntags: [imported]\n---\n")
	b.WriteString("# " + rd.Title + "\n\n")

	var meta []string
	if rd.PrepTime != "" {
		meta = append(meta, "**Prep Time:** "+rd.PrepTime)
	}
	if rd.CookTime != "" {
		meta = append(meta, "**Cook Time:** "+rd.CookTime)
	}
	if rd.Servings != "" {
		meta = append(meta, "**Servings:** "+rd.Servings)
	}
	if len(meta) > 0 {
		b.WriteString(strings.Join(meta, " | ") + "\n\n")
	}

	if len(rd.Ingredients) > 0 {
		b.WriteString("## Ingredients\n\n")
		for _, ing := range rd.Ingredients {
			b.WriteString("- " + ing + "\n")
		}
		b.WriteString("\n")
	}

	if len(rd.Instructions) > 0 {
		b.WriteString("## Instructions\n\n")
		for i, inst := range rd.Instructions {
			b.WriteString(fmt.Sprintf("%d. %s\n", i+1, inst))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Notes\n\n")
	b.WriteString(fmt.Sprintf("- Originally from [%s](%s)\n", sourceDomain, sourceURL))

	return b.String()
}
