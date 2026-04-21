package recipe

import "time"

type Recipe struct {
	Slug         string
	Title        string
	Category     string
	Tags         []string
	Image        string
	PrepTime     string
	CookTime     string
	Servings     string
	Ingredients  []Ingredient
	Instructions []string
	Notes        []string
	FilePath     string
	HTMLContent  string
	ModTime      time.Time
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
