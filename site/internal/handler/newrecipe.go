package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/whitenhiemer/recipe-site/internal/recipe"
	"github.com/whitenhiemer/recipe-site/internal/store"
)

var slugRegex = regexp.MustCompile(`[^a-z0-9]+`)

func toSlug(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = slugRegex.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func (s *Server) handleNewRecipe(w http.ResponseWriter, r *http.Request) {
	cats, _ := s.db.GetAllCategories()
	data := map[string]interface{}{
		"Title":      "New Recipe",
		"Categories": cats,
	}
	s.render(w, "new", data)
}

type newRecipeRequest struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Markdown string `json:"markdown"`
}

func (s *Server) handleCreateRecipeAPI(w http.ResponseWriter, r *http.Request) {
	var req newRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Category == "" || req.Markdown == "" {
		http.Error(w, "title, category, and markdown are required", http.StatusBadRequest)
		return
	}

	slug := toSlug(req.Title)
	if slug == "" {
		http.Error(w, "invalid title", http.StatusBadRequest)
		return
	}

	category := toSlug(req.Category)

	exists, err := s.db.RecipeExists(category, slug)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "a recipe with this name already exists in this category", http.StatusConflict)
		return
	}

	parsed, err := recipe.ParseMarkdown([]byte(req.Markdown), slug, category)
	if err != nil {
		http.Error(w, "failed to parse markdown", http.StatusBadRequest)
		return
	}

	now := time.Now()
	row := &store.RecipeRow{
		Slug:         slug,
		Title:        parsed.Title,
		Category:     category,
		Image:        parsed.Image,
		PrepTime:     parsed.PrepTime,
		CookTime:     parsed.CookTime,
		Servings:     parsed.Servings,
		Tags:         parsed.Tags,
		Ingredients:  toIngredientRows(parsed.Ingredients),
		Instructions: parsed.Instructions,
		Notes:        parsed.Notes,
		Markdown:     req.Markdown,
		HTMLContent:  parsed.HTMLContent,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if _, err := s.db.InsertRecipe(row); err != nil {
		http.Error(w, "failed to save recipe", http.StatusInternalServerError)
		return
	}

	if err := s.rebuildIndex(); err != nil {
		http.Error(w, "recipe saved but index rebuild failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"slug":     slug,
		"category": category,
	})
}

func toIngredientRows(ings []recipe.Ingredient) []store.IngredientRow {
	rows := make([]store.IngredientRow, len(ings))
	for i, ing := range ings {
		rows[i] = store.IngredientRow{
			Raw:      ing.Raw,
			Name:     ing.Name,
			Quantity: ing.Quantity,
			Unit:     ing.Unit,
		}
	}
	return rows
}

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

func (s *Server) handleImageUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "failed to read image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	ext, ok := allowedImageTypes[contentType]
	if !ok {
		http.Error(w, "only JPEG, PNG, and WebP images are allowed", http.StatusBadRequest)
		return
	}

	slug := toSlug(r.FormValue("slug"))
	if slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	filename := slug + ext
	imgDir := filepath.Join("static", "img", "recipes")
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		http.Error(w, "failed to create image directory", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		http.Error(w, "failed to read image", http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile(filepath.Join(imgDir, filename), buf.Bytes(), 0644); err != nil {
		http.Error(w, "failed to save image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"filename": filename,
		"url":      fmt.Sprintf("/static/img/recipes/%s", filename),
	})
}
