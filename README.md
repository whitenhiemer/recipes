# House Woodward Gourmand

Personal cooking recipe collection and self-hosted recipe website.

## Repository Structure

```
recipes/
├── breakfast/          # Morning meals and brunch
├── desserts/           # Sweets and baked goods
├── dinners/            # Main courses and entrees
├── doughs/             # Doughs and pastry bases
├── drinks/             # Cocktails, smoothies, and beverages
├── sauces-marinades/   # Sauces, dressings, rubs, brines, and marinades
├── sides/              # Side dishes and accompaniments
├── snacks/             # Appetizers and quick bites
├── soups-stews/        # Soups, stews, and chilis
├── TEMPLATE.md         # Recipe template
└── site/               # Go web application
    ├── main.go
    ├── Makefile
    ├── internal/       # Go packages (config, recipe parser, handlers)
    ├── templates/      # HTML templates (layout, home, recipe, etc.)
    ├── static/         # CSS, JS, images
    └── deploy/         # systemd, nginx, setup scripts
```

## Adding a Recipe

1. Copy [TEMPLATE.md](TEMPLATE.md) into the appropriate category folder
2. Rename it to match the recipe (e.g., `smoked-brisket.md`)
3. Fill in the sections following the template format
4. Commit and push -- the webhook will auto-reload the site

### Recipe Format

```markdown
---
tags: [tag1, tag2]
---
# Recipe Name

**Prep Time:** X min | **Cook Time:** X min | **Servings:** X

## Ingredients

- 1 cup ingredient

## Instructions

1. Step one

## Notes

- Tips, variations, substitutions
```

Tags are optional. Category is derived from the directory name. Subsection headers (`###`) can be used under Ingredients for multi-component recipes.

## Running the Site Locally

```bash
cd site
make run
# Opens at http://localhost:8080
```

## Deployment

See [site/deploy/](site/deploy/) for LXC deployment scripts and configuration.
See [PLAN.md](PLAN.md) for architecture details and project roadmap.
