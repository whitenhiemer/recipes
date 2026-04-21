package recipe

import (
	"sort"
	"strings"
	"sync"
)

type Index struct {
	mu      sync.RWMutex
	Recipes []*Recipe
	BySlug  map[string]*Recipe
	ByCat   map[string][]*Recipe
	ByTag   map[string][]*Recipe
	terms   map[string][]*Recipe
	AllTags []string
	AllCats []string
}

func NewIndex(recipes []*Recipe) *Index {
	idx := &Index{}
	idx.build(recipes)
	return idx
}

func (idx *Index) Rebuild(recipes []*Recipe) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.build(recipes)
}

func (idx *Index) build(recipes []*Recipe) {
	idx.Recipes = recipes
	idx.BySlug = make(map[string]*Recipe)
	idx.ByCat = make(map[string][]*Recipe)
	idx.ByTag = make(map[string][]*Recipe)
	idx.terms = make(map[string][]*Recipe)

	tagSet := make(map[string]bool)
	catSet := make(map[string]bool)

	for _, r := range recipes {
		idx.BySlug[r.Slug] = r

		if r.Category != "" {
			idx.ByCat[r.Category] = append(idx.ByCat[r.Category], r)
			catSet[r.Category] = true
		}

		for _, tag := range r.Tags {
			t := strings.ToLower(tag)
			idx.ByTag[t] = append(idx.ByTag[t], r)
			tagSet[t] = true
		}

		tokens := tokenize(r)
		seen := make(map[string]bool)
		for _, tok := range tokens {
			if !seen[tok] {
				idx.terms[tok] = append(idx.terms[tok], r)
				seen[tok] = true
			}
		}
	}

	idx.AllTags = sortedKeys(tagSet)
	idx.AllCats = sortedKeys(catSet)
}

func (idx *Index) Search(query string) []*Recipe {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	words := strings.Fields(strings.ToLower(query))
	if len(words) == 0 {
		return idx.Recipes
	}

	counts := make(map[string]int)
	for _, w := range words {
		for _, r := range idx.terms[w] {
			counts[r.Slug]++
		}
	}

	type scored struct {
		recipe *Recipe
		score  int
	}
	var results []scored
	for slug, count := range counts {
		results = append(results, scored{idx.BySlug[slug], count})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	out := make([]*Recipe, len(results))
	for i, s := range results {
		out[i] = s.recipe
	}
	return out
}

func (idx *Index) Reload(root string) error {
	recipes, err := ParseRecipesDir(root)
	if err != nil {
		return err
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.build(recipes)
	return nil
}

func (idx *Index) GetBySlug(slug string) *Recipe {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.BySlug[slug]
}

func (idx *Index) GetByCategory(cat string) []*Recipe {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.ByCat[cat]
}

func (idx *Index) GetByTag(tag string) []*Recipe {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.ByTag[strings.ToLower(tag)]
}

func (idx *Index) GetAllRecipes() []*Recipe {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.Recipes
}

func (idx *Index) GetAllTags() []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.AllTags
}

func (idx *Index) GetAllCategories() []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.AllCats
}

func (idx *Index) GenerateShoppingList(slugs []string) *ShoppingList {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	grouped := make(map[string][]string)
	order := []string{}

	for _, slug := range slugs {
		r, ok := idx.BySlug[slug]
		if !ok {
			continue
		}
		for _, ing := range r.Ingredients {
			name := strings.ToLower(ing.Name)
			if _, exists := grouped[name]; !exists {
				order = append(order, name)
			}
			amount := ing.Raw
			if r.Title != "" {
				amount = ing.Raw + " (" + r.Title + ")"
			}
			grouped[name] = append(grouped[name], amount)
		}
	}

	list := &ShoppingList{}
	for _, name := range order {
		list.Items = append(list.Items, ShoppingItem{
			Name:    name,
			Amounts: grouped[name],
		})
	}
	return list
}

func tokenize(r *Recipe) []string {
	var tokens []string

	tokens = append(tokens, strings.Fields(strings.ToLower(r.Title))...)
	tokens = append(tokens, strings.ToLower(r.Category))

	for _, tag := range r.Tags {
		tokens = append(tokens, strings.Fields(strings.ToLower(tag))...)
	}

	for _, ing := range r.Ingredients {
		tokens = append(tokens, strings.Fields(strings.ToLower(ing.Name))...)
	}

	return tokens
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
