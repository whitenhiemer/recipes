# House Woodward Gourmand -- Project Plan

## Overview

Self-hosted recipe website running on a Proxmox LXC container. Recipes are markdown files in this repo, parsed by a Go web app into a searchable, mobile-friendly site optimized for kitchen use.

## Architecture

- **Single repo** -- recipes at root, Go app in `site/`
- **No database** -- recipes parsed into in-memory inverted index on startup
- **Server-rendered HTML** with Go `html/template` + htmx for live search
- **Pico CSS** classless framework + custom mobile-first CSS
- **Meal plans in localStorage** -- single-user, no auth needed
- **Shopping lists generated server-side** from selected recipe slugs
- **GitHub webhook** auto-reloads index on push (no restart needed for recipe changes)

## Tech Stack

- Go 1.22+ (enhanced routing patterns)
- `github.com/yuin/goldmark` + `goldmark-meta` for markdown/frontmatter parsing
- htmx for live search
- Pico CSS for base styling
- Screen Wake Lock API for kitchen use

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
POST /webhook                       GitHub push webhook (HMAC-SHA256)
GET  /static/*                      Static assets
```

## LXC Deployment Target

- **Container:** Debian 12, 256MB RAM, 2GB disk, 1 CPU
- **Stack:** nginx (reverse proxy) -> Go binary on 127.0.0.1:8080
- **Service:** systemd with security hardening (NoNewPrivileges, ProtectSystem=strict)
- **TLS:** Let's Encrypt via certbot
- **Webhook secret:** stored in file at `/opt/recipe-site/.webhook-secret`
- **Layout:** single clone at `/opt/recipe-site/repo`, binary built from `repo/site/`

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
- [x] Meal plan page with localStorage (7-day grid, 3 meals/day)
- [x] Shopping list with checkboxes, print, copy-as-text
- [x] GitHub webhook handler with HMAC-SHA256 verification and rate limiting
- [x] Middleware: logging, recovery, security headers
- [x] systemd service, nginx config, deploy.sh, setup.sh
- [x] Branding: House Woodward Gourmand with family crest logo
- [x] Mobile-first responsive design (1/2/3 column breakpoints)
- [x] Touch-friendly UI (44-48px targets, large checkboxes)
- [x] Screen Wake Lock API on recipe pages (keep screen on while cooking)
- [x] 10 recipes across 7 categories

### TODO

- [ ] Deploy to Proxmox LXC container
- [ ] Configure GitHub webhook for auto-reload
- [ ] Set up TLS with Let's Encrypt
- [ ] Add unit tests for parser and index (`parser_test.go`, `index_test.go`)
- [ ] Add recipe images (photo field in frontmatter, display on cards and detail)
- [ ] Favorites (localStorage, filter view)
- [ ] "I made this" log/notes per recipe (localStorage)
- [ ] Prep timer / step-by-step cooking mode
- [ ] Dark/light theme toggle (Pico supports both)
- [ ] PWA manifest for "Add to Home Screen" on mobile
- [ ] Nutrition info parsing (optional frontmatter field)
- [ ] Print-optimized recipe detail page
- [ ] Recipe scaling (multiply ingredient quantities)
