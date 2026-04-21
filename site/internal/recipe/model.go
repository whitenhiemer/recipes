package recipe

import (
	"strings"
	"time"
)

type Recipe struct {
	Slug          string
	Title         string
	Category      string
	Tags          []string
	Image         string
	PrepTime      string
	CookTime      string
	Servings      string
	Ingredients   []Ingredient
	Instructions  []string
	Notes         []string
	FilePath      string
	HTMLContent   string
	SourceURL     string
	HasCapsaicin  bool
	ModTime       time.Time
}

var capsaicinKeywords = []string{
	"cayenne", "jalapeno", "jalapeño", "habanero", "serrano", "chipotle",
	"chili powder", "chili flakes", "chili pepper", "chilli",
	"red pepper flakes", "crushed red pepper", "hot sauce", "sriracha",
	"tabasco", "thai chili", "bird's eye", "ghost pepper", "scotch bonnet",
	"poblano", "ancho", "guajillo", "arbol", "pasilla", "hot pepper",
	"banana pepper", "cherry pepper", "fresno", "pepperoncini",
	"capsaicin", "paprika",
}

func DetectCapsaicin(ingredients []Ingredient) bool {
	for _, ing := range ingredients {
		name := strings.ToLower(ing.Name)
		raw := strings.ToLower(ing.Raw)
		for _, kw := range capsaicinKeywords {
			if strings.Contains(name, kw) || strings.Contains(raw, kw) {
				return true
			}
		}
	}
	return false
}

type Ingredient struct {
	Raw      string
	Name     string
	Quantity string
	Unit     string
}

type ShoppingList struct {
	Items []ShoppingItem
}

type ShoppingItem struct {
	Name    string
	Amounts []string
}
