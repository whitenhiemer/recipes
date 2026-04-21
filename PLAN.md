# House Woodward Gourmand -- Project Plan

## Overview

Self-hosted recipe website running on a Proxmox LXC container. Recipes stored in SQLite, with markdown import support. Go web app provides a searchable, mobile-friendly interface optimized for kitchen use.

## Architecture

- **Single repo** -- recipes at root (markdown), Go app in `site/`
- **SQLite database** -- primary storage for recipes (upgradeable to PostgreSQL)
- **Markdown import** -- on startup, markdown files are upserted into the DB
- **In-memory index** -- rebuilt from DB for fast search and filtering
- **Server-rendered HTML** with Go `html/template` + htmx for live search
- **Pico CSS** classless framework + custom mobile-first CSS
- **Meal plans in localStorage** -- single-user, no auth needed
- **Shopping lists generated server-side** from selected recipe slugs
- **GitHub webhook** triggers git pull + markdown reimport + index rebuild

## Tech Stack

- Go 1.22+ (enhanced routing patterns)
- `modernc.org/sqlite` -- pure-Go SQLite driver (no CGo needed)
- `github.com/yuin/goldmark` + `goldmark-meta` for markdown/frontmatter parsing
- htmx for live search
- Pico CSS for base styling
- Screen Wake Lock API for kitchen use

## Database

SQLite with WAL mode. Single `recipes` table storing all recipe data as JSON columns for ingredients, instructions, notes, and tags.

**Backup strategy:** Daily cron at 3 AM using `sqlite3 .backup`, 14-day retention. Backups stored in `/opt/recipe-site/data/backups/`.

**Migration path to PostgreSQL:** Replace `modernc.org/sqlite` with `lib/pq`, update connection string, adjust upsert syntax to use `ON CONFLICT ... DO UPDATE`. Schema is already PostgreSQL-compatible.

## Routes

```
GET  /                              Home: recipe grid, category tabs, search
GET  /recipes?q=&cat=&tag=          Filtered recipe list
GET  /recipes/{category}/{slug}     Recipe detail page
GET  /search?q=                     htmx partial: live search results
GET  /tags                          Tag listing
GET  /tags/{tag}                    Recipes by tag
GET  /mealplan                      Meal plan builder (JS + localStorage)
POST /api/shopping-list             JSON: {slugs: [...]} -> aggregated shopping list
GET  /shopping-list?slugs=a,b,c     Server-rendered shopping list
GET  /new                           Recipe creation wizard
POST /api/recipe                    Create recipe (JSON: title, category, markdown)
POST /api/image                     Upload recipe image
POST /webhook                       GitHub push webhook (HMAC-SHA256)
GET  /static/*                      Static assets
```

## LXC Deployment

### Automated (Proxmox community-scripts style)

Two scripts in `site/deploy/`:

- `ct-recipe-site.sh` -- runs on Proxmox host, creates LXC via whiptail dialogs
- `install-recipe-site.sh` -- runs inside container, installs everything

```bash
# From the Proxmox host
bash -c "$(cat /path/to/ct-recipe-site.sh)"
```

### Container Specs

- **OS:** Debian 12
- **RAM:** 512 MB
- **Disk:** 4 GB
- **CPU:** 1 core
- **Stack:** nginx -> Go binary on :8080 -> SQLite

### Layout

```
/opt/recipe-site/
├── recipe-site              # Go binary
├── .webhook-secret          # HMAC secret
├── repo/                    # Git clone of recipes repo
│   ├── site/                # Go source, templates, static assets
│   ├── breakfast/           # Markdown recipe files
│   └── ...
└── data/
    ├── recipes.db           # SQLite database
    └── backups/             # Daily .backup snapshots
```

## Progress

### Completed

- [x] Recipe markdown structure with YAML frontmatter tags
- [x] Go module, project structure under `site/`
- [x] Markdown parser with goldmark (title, metadata, ingredients, instructions, notes)
- [x] In-memory inverted index with search, category, and tag filtering
- [x] Shopping list generation (ingredient aggregation across recipes)
- [x] HTML templates: layout, home, recipe detail, search, tags, meal plan, shopping list
- [x] Live search with htmx
- [x] Category filter tabs
- [x] Tag pages
- [x] Meal plan page with localStorage (4 slots per meal, 3 categories)
- [x] Meal plan randomizer button
- [x] Shopping list with checkboxes, print, copy-as-text, save/load
- [x] Shopping list price estimates and emoji icons per ingredient
- [x] GitHub webhook handler with HMAC-SHA256 verification and rate limiting
- [x] Middleware: logging, recovery, security headers
- [x] systemd service, nginx config, deploy.sh
- [x] Branding: House Woodward Gourmand with family crest logo
- [x] Mobile-first responsive design (1/2/3 column breakpoints)
- [x] Touch-friendly UI (44-48px targets, large checkboxes)
- [x] Screen Wake Lock API on recipe pages
- [x] Recipe creation wizard (4-step form with preview)
- [x] Recipe image upload (JPEG/PNG/WebP, thumbnails on cards)
- [x] 29 recipes across 9 categories
- [x] SQLite database for recipe storage
- [x] Markdown file upload in wizard with syntax help guide
- [x] Proxmox community-scripts style deploy (ct + install scripts)
- [x] Database backup cron (daily, 14-day retention)

### TODO

- [ ] Deploy to Proxmox LXC container
- [ ] Configure GitHub webhook for auto-reload
- [ ] Set up TLS with Let's Encrypt
- [ ] Add unit tests for parser, index, and store
- [ ] Favorites (localStorage, filter view)
- [ ] "I made this" log/notes per recipe (localStorage)
- [ ] Prep timer / step-by-step cooking mode
- [ ] Dark/light theme toggle (Pico supports both)
- [ ] PWA manifest for "Add to Home Screen" on mobile
- [ ] Nutrition info parsing (optional frontmatter field)
- [ ] Print-optimized recipe detail page
- [ ] Recipe scaling (multiply ingredient quantities)
- [ ] Recipe editing/deletion from the web UI
- [ ] PostgreSQL migration (when needed)
