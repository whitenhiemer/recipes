package handler

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/whitenhiemer/recipe-site/internal/recipe"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Title":      "Recipes",
		"Recipes":    s.idx.GetAllRecipes(),
		"Categories": s.idx.GetAllCategories(),
		"Tags":       s.idx.GetAllTags(),
		"ActiveCat":  "",
	}
	s.render(w, "home", data)
}

func (s *Server) handleRecipeList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	cat := r.URL.Query().Get("cat")
	tag := r.URL.Query().Get("tag")

	var recipes = s.idx.GetAllRecipes()

	if q != "" {
		recipes = s.idx.Search(q)
	}
	if cat != "" {
		recipes = s.idx.GetByCategory(cat)
	}
	if tag != "" {
		recipes = s.idx.GetByTag(tag)
	}

	data := map[string]interface{}{
		"Title":      "Recipes",
		"Recipes":    recipes,
		"Categories": s.idx.GetAllCategories(),
		"Tags":       s.idx.GetAllTags(),
		"Query":      q,
		"ActiveCat":  cat,
		"ActiveTag":  tag,
	}
	s.render(w, "home", data)
}

func (s *Server) handleRecipeDetail(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	recipe := s.idx.GetBySlug(slug)
	if recipe == nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Title":   recipe.Title,
		"Recipe":  recipe,
		"Content": template.HTML(recipe.HTMLContent),
	}
	s.render(w, "recipe", data)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	results := s.idx.Search(q)

	data := map[string]interface{}{
		"Recipes": results,
		"Query":   q,
	}
	tmpl := s.templates["search"]
	if tmpl == nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "search.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleTags(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Tags",
		"Tags":  s.idx.GetAllTags(),
	}
	s.render(w, "tags", data)
}

func (s *Server) handleTagRecipes(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	recipes := s.idx.GetByTag(tag)

	data := map[string]interface{}{
		"Title":   "Tag: " + tag,
		"Recipes": recipes,
		"Tags":    s.idx.GetAllTags(),
		"Tag":     tag,
	}
	s.render(w, "home", data)
}

func (s *Server) handleMealPlan(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":            "Meal Plan",
		"Recipes":          s.idx.GetAllRecipes(),
		"BreakfastRecipes": s.idx.GetByCategory("breakfast"),
		"LunchRecipes":     s.idx.GetByCategory("lunch"),
		"DinnerRecipes":    s.idx.GetByCategory("dinners"),
	}
	s.render(w, "mealplan", data)
}

func (s *Server) handlePantry(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Pantry & Fridge",
	}
	s.render(w, "pantry", data)
}

func (s *Server) handleShoppingList(w http.ResponseWriter, r *http.Request) {
	slugParam := r.URL.Query().Get("slugs")
	trip1Param := r.URL.Query().Get("trip1")
	trip2Param := r.URL.Query().Get("trip2")

	hasTrips := trip1Param != "" || trip2Param != ""

	if slugParam == "" && !hasTrips {
		data := map[string]interface{}{
			"Title":    "Shopping List",
			"Items":    nil,
			"HasSaved": true,
		}
		s.render(w, "shopping", data)
		return
	}

	if hasTrips {
		var t1Slugs, t2Slugs []string
		if trip1Param != "" {
			t1Slugs = strings.Split(trip1Param, ",")
		}
		if trip2Param != "" {
			t2Slugs = strings.Split(trip2Param, ",")
		}

		trip1List, trip2List := s.idx.GenerateShoppingTrips(t1Slugs, t2Slugs)
		trip1Priced, trip1Total := recipe.PriceShoppingList(trip1List)
		trip2Priced, trip2Total := recipe.PriceShoppingList(trip2List)
		trip1Depts := recipe.GroupByDepartment(trip1Priced)
		trip2Depts := recipe.GroupByDepartment(trip2Priced)

		data := map[string]interface{}{
			"Title":           "Shopping List",
			"Items":           trip1List.Items,
			"HasTrips":        true,
			"Trip1Departments": trip1Depts,
			"Trip1Total":      trip1Total,
			"Trip2Departments": trip2Depts,
			"Trip2Total":      trip2Total,
			"EstimatedTotal":  trip1Total + trip2Total,
			"HasSaved":        false,
		}
		s.render(w, "shopping", data)
		return
	}

	slugs := strings.Split(slugParam, ",")
	list := s.idx.GenerateShoppingList(slugs)
	pricedItems, total := recipe.PriceShoppingList(list)
	departments := recipe.GroupByDepartment(pricedItems)

	data := map[string]interface{}{
		"Title":          "Shopping List",
		"Items":          list.Items,
		"PricedItems":    pricedItems,
		"Departments":    departments,
		"EstimatedTotal": total,
		"HasSaved":       false,
	}
	s.render(w, "shopping", data)
}
