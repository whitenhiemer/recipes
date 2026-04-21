package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var slugRegex = regexp.MustCompile(`[^a-z0-9]+`)

func toSlug(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = slugRegex.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func (s *Server) handleNewRecipe(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":      "New Recipe",
		"Categories": s.idx.GetAllCategories(),
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
	dir := filepath.Join(s.cfg.RecipesDir, category)

	if err := os.MkdirAll(dir, 0755); err != nil {
		http.Error(w, "failed to create category directory", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(dir, slug+".md")

	if _, err := os.Stat(filePath); err == nil {
		http.Error(w, "a recipe with this name already exists in this category", http.StatusConflict)
		return
	}

	if err := os.WriteFile(filePath, []byte(req.Markdown), 0644); err != nil {
		http.Error(w, "failed to write recipe file", http.StatusInternalServerError)
		return
	}

	if err := s.idx.Reload(s.cfg.RecipesDir); err != nil {
		http.Error(w, "recipe saved but index reload failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"slug":     slug,
		"category": category,
	})
}

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

func (s *Server) handleImageUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10MB max

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

	dst, err := os.Create(filepath.Join(imgDir, filename))
	if err != nil {
		http.Error(w, "failed to save image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "failed to write image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"filename": filename,
		"url":      fmt.Sprintf("/static/img/recipes/%s", filename),
	})
}
