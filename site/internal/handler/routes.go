package handler

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/whitenhiemer/recipe-site/internal/config"
	"github.com/whitenhiemer/recipe-site/internal/recipe"
	"github.com/whitenhiemer/recipe-site/internal/store"
)

type Server struct {
	idx       *recipe.Index
	cfg       *config.Config
	db        *store.DB
	templates map[string]*template.Template
}

func Register(mux *http.ServeMux, idx *recipe.Index, cfg *config.Config, db *store.DB) {
	s := &Server{idx: idx, cfg: cfg, db: db}
	s.loadTemplates()

	mux.HandleFunc("GET /", s.handleHome)
	mux.HandleFunc("GET /recipes", s.handleRecipeList)
	mux.HandleFunc("GET /recipes/{category}/{slug}", s.handleRecipeDetail)
	mux.HandleFunc("GET /search", s.handleSearch)
	mux.HandleFunc("GET /tags", s.handleTags)
	mux.HandleFunc("GET /tags/{tag}", s.handleTagRecipes)
	mux.HandleFunc("GET /mealplan", s.handleMealPlan)
	mux.HandleFunc("POST /api/shopping-list", s.handleShoppingListAPI)
	mux.HandleFunc("GET /shopping-list", s.handleShoppingList)
	mux.HandleFunc("GET /new", s.handleNewRecipe)
	mux.HandleFunc("POST /api/recipe", s.handleCreateRecipeAPI)
	mux.HandleFunc("POST /api/image", s.handleImageUpload)
	mux.HandleFunc("POST /webhook", s.handleWebhook)

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func (s *Server) loadTemplates() {
	layout := filepath.Join("templates", "layout.html")
	partials, _ := filepath.Glob(filepath.Join("templates", "partials", "*.html"))

	base := append([]string{layout}, partials...)

	pages := []string{"home", "recipe", "search", "tags", "mealplan", "shopping", "new"}
	s.templates = make(map[string]*template.Template)

	for _, page := range pages {
		files := make([]string, len(base))
		copy(files, base)
		files = append(files, filepath.Join("templates", page+".html"))
		s.templates[page] = template.Must(template.ParseFiles(files...))
	}
}

func (s *Server) render(w http.ResponseWriter, page string, data any) {
	tmpl, ok := s.templates[page]
	if !ok {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) rebuildIndex() error {
	rows, err := s.db.GetAllRecipes()
	if err != nil {
		return err
	}
	s.idx.Rebuild(store.RowsToRecipes(rows))
	log.Printf("index rebuilt: %d recipes", len(rows))
	return nil
}
