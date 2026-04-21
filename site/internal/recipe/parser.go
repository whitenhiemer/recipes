package recipe

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var (
	metaRegex = regexp.MustCompile(`(?:\*\*)?Prep Time:(?:\*\*)?\s*(.+?)\s*\|\s*(?:\*\*)?Cook Time:(?:\*\*)?\s*(.+?)\s*\|\s*(?:\*\*)?Servings:(?:\*\*)?\s*(.+)`)
	unitRegex = regexp.MustCompile(`(?i)^([\d/\.\s½¼¾⅓⅔]+)?\s*(cups?|tbsp|tsp|oz|lbs?|g|kg|ml|cloves?|cans?|bunch|pinch|dash|slices?|pieces?|stalks?|heads?)?\s*[,.]?\s*(.+)$`)
)

func ParseRecipesDir(root string) ([]*Recipe, error) {
	var recipes []*Recipe

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			rel, _ := filepath.Rel(root, path)
			if rel == "site" || strings.HasPrefix(rel, "site/") || rel == ".git" || strings.HasPrefix(rel, ".git/") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		if filepath.Dir(rel) == "." {
			return nil
		}

		r, err := ParseRecipeFile(path, root)
		if err != nil {
			return nil
		}
		recipes = append(recipes, r)
		return nil
	})

	return recipes, err
}

func ParseRecipeFile(path string, root string) (*Recipe, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	rel, _ := filepath.Rel(root, path)
	dir := filepath.Dir(rel)
	base := filepath.Base(rel)
	slug := strings.TrimSuffix(base, ".md")

	category := ""
	if dir != "." {
		category = dir
	}

	r, err := ParseMarkdown(data, slug, category)
	if err != nil {
		return nil, err
	}
	r.FilePath = rel
	r.ModTime = info.ModTime()
	return r, nil
}

func ParseMarkdown(data []byte, slug, category string) (*Recipe, error) {
	md := goldmark.New(goldmark.WithExtensions(meta.Meta))
	ctx := parser.NewContext()

	var htmlBuf bytes.Buffer
	if err := md.Convert(data, &htmlBuf, parser.WithContext(ctx)); err != nil {
		return nil, err
	}

	r := &Recipe{
		Slug:        slug,
		Category:    category,
		HTMLContent: htmlBuf.String(),
	}

	metaData := meta.Get(ctx)
	if tags, ok := metaData["tags"]; ok {
		switch v := tags.(type) {
		case []interface{}:
			for _, t := range v {
				if s, ok := t.(string); ok {
					r.Tags = append(r.Tags, s)
				}
			}
		}
	}
	if img, ok := metaData["image"]; ok {
		if s, ok := img.(string); ok {
			r.Image = s
		}
	}
	if r.Tags == nil {
		r.Tags = []string{}
	}

	reader := text.NewReader(data)
	doc := md.Parser().Parse(reader, parser.WithContext(parser.NewContext()))
	parseAST(doc, data, r)

	if r.Instructions == nil {
		r.Instructions = []string{}
	}
	if r.Notes == nil {
		r.Notes = []string{}
	}

	r.HasCapsaicin = DetectCapsaicin(r.Ingredients)

	return r, nil
}

func parseAST(doc ast.Node, source []byte, r *Recipe) {
	var currentSection string

	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Heading:
			if node.Level == 1 && r.Title == "" {
				r.Title = extractText(node, source)
			}
			if node.Level == 2 {
				currentSection = strings.ToLower(extractText(node, source))
			}

		case *ast.Paragraph:
			if currentSection == "" {
				text := extractText(node, source)
				if m := metaRegex.FindStringSubmatch(text); m != nil {
					r.PrepTime = strings.TrimSpace(m[1])
					r.CookTime = strings.TrimSpace(m[2])
					r.Servings = strings.TrimSpace(m[3])
				}
			}

		case *ast.ListItem:
			text := extractText(node, source)
			switch {
			case strings.HasPrefix(currentSection, "ingredient"):
				ing := parseIngredient(text)
				r.Ingredients = append(r.Ingredients, ing)
			case strings.HasPrefix(currentSection, "instruction"):
				r.Instructions = append(r.Instructions, text)
			case strings.HasPrefix(currentSection, "note"):
				r.Notes = append(r.Notes, text)
			}
		}

		return ast.WalkContinue, nil
	})
}

func extractText(node ast.Node, source []byte) string {
	var buf bytes.Buffer
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if t, ok := child.(*ast.Text); ok {
			buf.Write(t.Segment.Value(source))
		} else if _, ok := child.(*ast.Emphasis); ok {
			buf.WriteString(extractText(child, source))
		} else if _, ok := child.(*ast.CodeSpan); ok {
			buf.WriteString(extractText(child, source))
		} else {
			buf.WriteString(extractText(child, source))
		}
	}
	return strings.TrimSpace(buf.String())
}

func parseIngredient(raw string) Ingredient {
	ing := Ingredient{Raw: raw}

	m := unitRegex.FindStringSubmatch(raw)
	if m != nil {
		ing.Quantity = strings.TrimSpace(m[1])
		ing.Unit = strings.TrimSpace(m[2])
		ing.Name = strings.TrimSpace(m[3])
	} else {
		ing.Name = raw
	}

	return ing
}
