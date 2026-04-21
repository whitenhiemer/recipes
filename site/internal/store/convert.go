package store

import (
	"github.com/whitenhiemer/recipe-site/internal/recipe"
)

func RowsToRecipes(rows []*RecipeRow) []*recipe.Recipe {
	recipes := make([]*recipe.Recipe, len(rows))
	for i, row := range rows {
		recipes[i] = RowToRecipe(row)
	}
	return recipes
}

func RowToRecipe(row *RecipeRow) *recipe.Recipe {
	ings := make([]recipe.Ingredient, len(row.Ingredients))
	for i, ing := range row.Ingredients {
		ings[i] = recipe.Ingredient{
			Raw:      ing.Raw,
			Name:     ing.Name,
			Quantity: ing.Quantity,
			Unit:     ing.Unit,
		}
	}

	r := &recipe.Recipe{
		Slug:         row.Slug,
		Title:        row.Title,
		Category:     row.Category,
		Tags:         row.Tags,
		Image:        row.Image,
		PrepTime:     row.PrepTime,
		CookTime:     row.CookTime,
		Servings:     row.Servings,
		Ingredients:  ings,
		Instructions: row.Instructions,
		Notes:        row.Notes,
		HTMLContent:  row.HTMLContent,
		SourceURL:    row.SourceURL,
		ModTime:      row.UpdatedAt,
	}
	r.HasCapsaicin = recipe.DetectCapsaicin(ings)
	return r
}
