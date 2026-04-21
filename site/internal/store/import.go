package store

import (
	"log"
	"time"

	"github.com/whitenhiemer/recipe-site/internal/recipe"
)

func (d *DB) ImportFromMarkdown(recipesDir string) (int, error) {
	recipes, err := recipe.ParseRecipesDir(recipesDir)
	if err != nil {
		return 0, err
	}

	var imported int
	for _, r := range recipes {
		row := recipeToRow(r)
		err := d.UpsertRecipe(row)
		if err != nil {
			log.Printf("import: failed to upsert %s/%s: %v", r.Category, r.Slug, err)
			continue
		}
		imported++
	}

	return imported, nil
}

func recipeToRow(r *recipe.Recipe) *RecipeRow {
	ings := make([]IngredientRow, len(r.Ingredients))
	for i, ing := range r.Ingredients {
		ings[i] = IngredientRow{
			Raw:      ing.Raw,
			Name:     ing.Name,
			Quantity: ing.Quantity,
			Unit:     ing.Unit,
		}
	}

	now := time.Now()
	modTime := r.ModTime
	if modTime.IsZero() {
		modTime = now
	}

	return &RecipeRow{
		Slug:         r.Slug,
		Title:        r.Title,
		Category:     r.Category,
		Image:        r.Image,
		PrepTime:     r.PrepTime,
		CookTime:     r.CookTime,
		Servings:     r.Servings,
		Tags:         r.Tags,
		Ingredients:  ings,
		Instructions: r.Instructions,
		Notes:        r.Notes,
		Markdown:     "",
		HTMLContent:  r.HTMLContent,
		CreatedAt:    modTime,
		UpdatedAt:    modTime,
	}
}
