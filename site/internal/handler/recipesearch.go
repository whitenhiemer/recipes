package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type searchResult struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Desc  string `json:"desc"`
}

var (
	ddgResultPattern = regexp.MustCompile(`<a[^>]+class="result__a"[^>]+href="([^"]+)"[^>]*>(.*?)</a>`)
	ddgSnippetPattern = regexp.MustCompile(`<a[^>]+class="result__snippet"[^>]*>(.*?)</a>`)
	htmlTagPattern    = regexp.MustCompile(`<[^>]*>`)
)

func (s *Server) handleRecipeSearch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	query := strings.TrimSpace(req.Query)
	if query == "" {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	searchQuery := query + " recipe"
	ddgURL := "https://html.duckduckgo.com/html/?q=" + url.QueryEscape(searchQuery)

	client := &http.Client{Timeout: 10 * time.Second}
	httpReq, _ := http.NewRequest("GET", ddgURL, nil)
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, "search failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		http.Error(w, "failed to read search results", http.StatusBadGateway)
		return
	}

	html := string(body)
	results := parseSearchResults(html, 3)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func parseSearchResults(html string, limit int) []searchResult {
	linkMatches := ddgResultPattern.FindAllStringSubmatch(html, -1)
	snippetMatches := ddgSnippetPattern.FindAllStringSubmatch(html, -1)

	var results []searchResult
	for i, m := range linkMatches {
		if len(results) >= limit {
			break
		}

		rawURL := m[1]
		realURL := extractDDGRedirectURL(rawURL)
		if realURL == "" {
			continue
		}

		if !looksLikeRecipeSite(realURL) {
			continue
		}

		title := htmlTagPattern.ReplaceAllString(m[2], "")
		title = strings.TrimSpace(title)

		desc := ""
		if i < len(snippetMatches) {
			desc = htmlTagPattern.ReplaceAllString(snippetMatches[i][1], "")
			desc = strings.TrimSpace(desc)
		}

		results = append(results, searchResult{
			Title: title,
			URL:   realURL,
			Desc:  desc,
		})
	}

	return results
}

func extractDDGRedirectURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	if uddg := parsed.Query().Get("uddg"); uddg != "" {
		return uddg
	}

	if parsed.Scheme == "http" || parsed.Scheme == "https" {
		return rawURL
	}

	return ""
}

func looksLikeRecipeSite(u string) bool {
	recipeDomains := []string{
		"allrecipes.com", "foodnetwork.com", "epicurious.com",
		"bonappetit.com", "seriouseats.com", "simplyrecipes.com",
		"tasty.co", "delish.com", "food52.com", "budgetbytes.com",
		"cookieandkate.com", "minimalistbaker.com", "skinnytaste.com",
		"thekitchn.com", "marthastewart.com", "kingarthurbaking.com",
		"bettycrocker.com", "pillsbury.com", "myrecipes.com",
		"eatingwell.com", "cooking.nytimes.com", "bbcgoodfood.com",
		"taste.com.au", "recipetineats.com", "damndelicious.net",
		"pinchofyum.com", "halfbakedharvest.com", "smittenkitchen.com",
		"loveandlemons.com", "cafedelites.com", "sallysbakingaddiction.com",
	}

	lower := strings.ToLower(u)
	for _, domain := range recipeDomains {
		if strings.Contains(lower, domain) {
			return true
		}
	}

	recipeKeywords := []string{"recipe", "cooking", "food", "baking", "kitchen", "cook", "chef"}
	for _, kw := range recipeKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}

	parsed, err := url.Parse(u)
	if err != nil {
		return false
	}
	return parsed.Scheme == "https" && !strings.Contains(lower, "youtube.com") && !strings.Contains(lower, "amazon.com") && !strings.Contains(lower, "wikipedia.org")
}
