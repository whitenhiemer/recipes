package store

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS recipes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			slug TEXT NOT NULL,
			title TEXT NOT NULL,
			category TEXT NOT NULL,
			image TEXT NOT NULL DEFAULT '',
			prep_time TEXT NOT NULL DEFAULT '',
			cook_time TEXT NOT NULL DEFAULT '',
			servings TEXT NOT NULL DEFAULT '',
			tags TEXT NOT NULL DEFAULT '[]',
			ingredients TEXT NOT NULL DEFAULT '[]',
			instructions TEXT NOT NULL DEFAULT '[]',
			notes TEXT NOT NULL DEFAULT '[]',
			markdown TEXT NOT NULL DEFAULT '',
			html_content TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(category, slug)
		);
		CREATE INDEX IF NOT EXISTS idx_recipes_slug ON recipes(slug);
		CREATE INDEX IF NOT EXISTS idx_recipes_category ON recipes(category);
	`)
	return err
}

type RecipeRow struct {
	ID           int64
	Slug         string
	Title        string
	Category     string
	Image        string
	PrepTime     string
	CookTime     string
	Servings     string
	Tags         []string
	Ingredients  []IngredientRow
	Instructions []string
	Notes        []string
	Markdown     string
	HTMLContent  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type IngredientRow struct {
	Raw      string `json:"raw"`
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	Unit     string `json:"unit"`
}

func (d *DB) InsertRecipe(r *RecipeRow) (int64, error) {
	tagsJSON, _ := json.Marshal(r.Tags)
	ingsJSON, _ := json.Marshal(r.Ingredients)
	instJSON, _ := json.Marshal(r.Instructions)
	notesJSON, _ := json.Marshal(r.Notes)

	result, err := d.db.Exec(`
		INSERT INTO recipes (slug, title, category, image, prep_time, cook_time, servings,
			tags, ingredients, instructions, notes, markdown, html_content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.Slug, r.Title, r.Category, r.Image, r.PrepTime, r.CookTime, r.Servings,
		string(tagsJSON), string(ingsJSON), string(instJSON), string(notesJSON),
		r.Markdown, r.HTMLContent, r.CreatedAt, r.UpdatedAt,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (d *DB) UpdateRecipe(r *RecipeRow) error {
	tagsJSON, _ := json.Marshal(r.Tags)
	ingsJSON, _ := json.Marshal(r.Ingredients)
	instJSON, _ := json.Marshal(r.Instructions)
	notesJSON, _ := json.Marshal(r.Notes)

	_, err := d.db.Exec(`
		UPDATE recipes SET title=?, image=?, prep_time=?, cook_time=?, servings=?,
			tags=?, ingredients=?, instructions=?, notes=?, markdown=?, html_content=?,
			updated_at=?
		WHERE category=? AND slug=?`,
		r.Title, r.Image, r.PrepTime, r.CookTime, r.Servings,
		string(tagsJSON), string(ingsJSON), string(instJSON), string(notesJSON),
		r.Markdown, r.HTMLContent, r.UpdatedAt,
		r.Category, r.Slug,
	)
	return err
}

func (d *DB) UpsertRecipe(r *RecipeRow) error {
	tagsJSON, _ := json.Marshal(r.Tags)
	ingsJSON, _ := json.Marshal(r.Ingredients)
	instJSON, _ := json.Marshal(r.Instructions)
	notesJSON, _ := json.Marshal(r.Notes)

	_, err := d.db.Exec(`
		INSERT INTO recipes (slug, title, category, image, prep_time, cook_time, servings,
			tags, ingredients, instructions, notes, markdown, html_content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(category, slug) DO UPDATE SET
			title=excluded.title, image=excluded.image,
			prep_time=excluded.prep_time, cook_time=excluded.cook_time, servings=excluded.servings,
			tags=excluded.tags, ingredients=excluded.ingredients,
			instructions=excluded.instructions, notes=excluded.notes,
			markdown=excluded.markdown, html_content=excluded.html_content,
			updated_at=excluded.updated_at`,
		r.Slug, r.Title, r.Category, r.Image, r.PrepTime, r.CookTime, r.Servings,
		string(tagsJSON), string(ingsJSON), string(instJSON), string(notesJSON),
		r.Markdown, r.HTMLContent, r.CreatedAt, r.UpdatedAt,
	)
	return err
}

func (d *DB) GetAllRecipes() ([]*RecipeRow, error) {
	return d.queryRecipes("SELECT * FROM recipes ORDER BY category, title")
}

func (d *DB) GetByCategory(cat string) ([]*RecipeRow, error) {
	return d.queryRecipes("SELECT * FROM recipes WHERE category=? ORDER BY title", cat)
}

func (d *DB) GetBySlug(slug string) (*RecipeRow, error) {
	rows, err := d.queryRecipes("SELECT * FROM recipes WHERE slug=? LIMIT 1", slug)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return rows[0], nil
}

func (d *DB) GetByCategoryAndSlug(cat, slug string) (*RecipeRow, error) {
	rows, err := d.queryRecipes("SELECT * FROM recipes WHERE category=? AND slug=? LIMIT 1", cat, slug)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return rows[0], nil
}

func (d *DB) DeleteRecipe(category, slug string) error {
	_, err := d.db.Exec("DELETE FROM recipes WHERE category=? AND slug=?", category, slug)
	return err
}

func (d *DB) SearchRecipes(query string) ([]*RecipeRow, error) {
	words := strings.Fields(strings.ToLower(query))
	if len(words) == 0 {
		return d.GetAllRecipes()
	}

	conditions := make([]string, len(words))
	args := make([]interface{}, len(words)*4)
	for i, w := range words {
		pattern := "%" + w + "%"
		conditions[i] = "(LOWER(title) LIKE ? OR LOWER(category) LIKE ? OR LOWER(tags) LIKE ? OR LOWER(ingredients) LIKE ?)"
		args[i*4] = pattern
		args[i*4+1] = pattern
		args[i*4+2] = pattern
		args[i*4+3] = pattern
	}

	q := "SELECT * FROM recipes WHERE " + strings.Join(conditions, " AND ") + " ORDER BY title"
	return d.queryRecipes(q, args...)
}

func (d *DB) GetAllTags() ([]string, error) {
	rows, err := d.db.Query("SELECT DISTINCT tags FROM recipes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagSet := make(map[string]bool)
	for rows.Next() {
		var tagsJSON string
		if err := rows.Scan(&tagsJSON); err != nil {
			continue
		}
		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}
		for _, t := range tags {
			tagSet[strings.ToLower(t)] = true
		}
	}

	result := make([]string, 0, len(tagSet))
	for t := range tagSet {
		if t != "" {
			result = append(result, t)
		}
	}
	return result, nil
}

func (d *DB) GetAllCategories() ([]string, error) {
	rows, err := d.db.Query("SELECT DISTINCT category FROM recipes ORDER BY category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			continue
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

func (d *DB) GetByTag(tag string) ([]*RecipeRow, error) {
	pattern := `%"` + strings.ToLower(tag) + `"%`
	return d.queryRecipes("SELECT * FROM recipes WHERE LOWER(tags) LIKE ? ORDER BY title", pattern)
}

func (d *DB) RecipeExists(category, slug string) (bool, error) {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM recipes WHERE category=? AND slug=?", category, slug).Scan(&count)
	return count > 0, err
}

func (d *DB) queryRecipes(query string, args ...interface{}) ([]*RecipeRow, error) {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*RecipeRow
	for rows.Next() {
		r := &RecipeRow{}
		var tagsJSON, ingsJSON, instJSON, notesJSON string
		err := rows.Scan(
			&r.ID, &r.Slug, &r.Title, &r.Category, &r.Image,
			&r.PrepTime, &r.CookTime, &r.Servings,
			&tagsJSON, &ingsJSON, &instJSON, &notesJSON,
			&r.Markdown, &r.HTMLContent,
			&r.CreatedAt, &r.UpdatedAt,
		)
		if err != nil {
			continue
		}
		json.Unmarshal([]byte(tagsJSON), &r.Tags)
		json.Unmarshal([]byte(ingsJSON), &r.Ingredients)
		json.Unmarshal([]byte(instJSON), &r.Instructions)
		json.Unmarshal([]byte(notesJSON), &r.Notes)
		if r.Tags == nil {
			r.Tags = []string{}
		}
		if r.Ingredients == nil {
			r.Ingredients = []IngredientRow{}
		}
		if r.Instructions == nil {
			r.Instructions = []string{}
		}
		if r.Notes == nil {
			r.Notes = []string{}
		}
		recipes = append(recipes, r)
	}
	return recipes, nil
}
