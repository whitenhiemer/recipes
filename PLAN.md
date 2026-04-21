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
GET  /pantry                        Pantry & fridge inventory (localStorage)
POST /api/recipe                    Create recipe (JSON: title, category, markdown)
POST /api/recipe/import-url         Import recipe from external URL (JSON-LD extraction)
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
- [x] Shopping list grouped by grocery store departments
- [x] Pantry staples designation (salt, oil, flour, spices auto-excluded from total)
- [x] Branded print layout with logo (single-page, 3-column, compact)
- [x] Pantry & fridge inventory page (localStorage, add/remove/filter by location)
- [x] On-hand item detection in shopping list (fuzzy match against pantry inventory)
- [x] Buy units -- standard purchase quantities per ingredient (e.g., "1 lb / 4 sticks")
- [x] Prepopulated pantry suggestions with emoji icons, grouped by fridge/pantry/freezer
- [x] Two-trip shopping lists (shelf-stable in trip 1, perishables split by early/late week)
- [x] Gas tank reminder for Libby on shopping list printouts
- [x] Meal plan dropdown selects with category-grouped recipe choices
- [x] Recipe URL import via JSON-LD extraction with source attribution
- [x] Google recipe search integration on import tab
- [x] Recipe scaling (default 5 servings, +/- control, fraction display)
- [x] Fix metadata regex to parse Servings/PrepTime/CookTime from AST text

### TODO

- [ ] Deploy to Proxmox LXC container
- [ ] Configure GitHub webhook for auto-reload
- [ ] Set up TLS with Let's Encrypt
- [ ] Acquire a domain name for the site
- [ ] Google OAuth2 authentication (see Auth section below)
- [ ] Add unit tests for parser, index, and store
- [ ] Favorites (localStorage, filter view)
- [ ] "I made this" log/notes per recipe (localStorage)
- [ ] Prep timer / step-by-step cooking mode
- [ ] Dark/light theme toggle (Pico supports both)
- [ ] PWA manifest for "Add to Home Screen" on mobile
- [ ] Nutrition info parsing (optional frontmatter field)
- [ ] Print-optimized recipe detail page
- [ ] Recipe editing/deletion from the web UI
- [ ] PostgreSQL migration (when needed)
- [ ] Share recipes via SMS/WhatsApp (Web Share API or Twilio)
- [ ] Compose recipes from existing sub-recipes (e.g., reuse dough/sauce recipes as ingredients in new recipes)
- [ ] AI recipe generation from an idea or prompt (e.g., "something with chicken and mushrooms")

## Authentication (Google OAuth2)

**Goal:** Public read access to all recipes. Write access (create, edit, delete) requires Google login from an allowlisted email.

### Prerequisites

- Domain name with TLS (Let's Encrypt) -- required for OAuth redirect URI
- Google Cloud project with OAuth 2.0 credentials (client ID + secret)

### Google Cloud Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a project (or use existing): "House Woodward Gourmand"
3. APIs & Services > OAuth consent screen > External > configure app name, support email
4. APIs & Services > Credentials > Create OAuth 2.0 Client ID
   - Application type: Web application
   - Authorized redirect URI: `https://recipes.yourdomain.com/auth/callback`
5. Save the Client ID and Client Secret

### New Routes

```
GET  /login              Redirect to Google OAuth consent screen
GET  /auth/callback      Exchange auth code for token, create session, redirect to /
GET  /logout             Clear session cookie, redirect to /
```

### New Package: `internal/auth/`

```
internal/auth/
├── oauth.go             # Google OAuth2 config, login/callback/logout handlers
├── session.go           # Session create/get/delete, cookie management
└── middleware.go         # RequireAuth middleware for write routes
```

### Database Changes

New `sessions` table:

```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,           -- random 32-byte hex token
    email TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL
);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
```

Session cleanup: delete expired rows on each login or via periodic goroutine.

### Config Changes

New flags/env vars:

```
-google-client-id       or  GOOGLE_CLIENT_ID
-google-client-secret   or  GOOGLE_CLIENT_SECRET
-allowed-emails         or  ALLOWED_EMAILS       (comma-separated)
-base-url               or  BASE_URL             (e.g. https://recipes.yourdomain.com)
```

Store OAuth client secret in a file (like webhook secret) to avoid CLI args in `ps` output:

```
-google-secret-file=/opt/recipe-site/.google-oauth-secret
```

### Auth Flow

1. User clicks "Login" in nav (shown when not authenticated)
2. `GET /login` redirects to Google with `state` parameter (CSRF token stored in cookie)
3. Google redirects to `/auth/callback?code=...&state=...`
4. Server exchanges code for token, fetches user info (email)
5. If email is in allowlist: create session row in DB, set secure cookie, redirect to /
6. If email is not in allowlist: show "not authorized" page
7. Session cookie: `HttpOnly`, `Secure`, `SameSite=Lax`, 7-day expiry

### Protected Routes

Only these routes require auth -- everything else stays public:

```
GET  /new                RequireAuth
POST /api/recipe         RequireAuth
POST /api/image          RequireAuth
POST /api/recipe/delete  RequireAuth (future)
POST /api/recipe/edit    RequireAuth (future)
```

### UI Changes

- Nav: show "Login" link when no session, show email + "Logout" when authenticated
- "+ New Recipe" button: still visible to everyone, but clicking it redirects to /login if not authenticated
- Recipe detail: show "Edit" / "Delete" buttons only when authenticated (future)

### Dependencies

```
golang.org/x/oauth2        # OAuth2 client
```

No other external deps needed. Google's userinfo endpoint returns email directly.

### Deploy Changes

- Store Google OAuth credentials in `/opt/recipe-site/.google-oauth-secret`
- Add `-google-client-id`, `-google-secret-file`, `-allowed-emails`, `-base-url` to systemd ExecStart
- nginx must proxy the auth routes (already covered by `location /`)

### Security Notes

- CSRF protection via `state` parameter in OAuth flow (random token in cookie, verified on callback)
- Session tokens are cryptographically random (32 bytes from crypto/rand)
- Cookies are HttpOnly + Secure + SameSite=Lax
- Email allowlist is checked server-side on every callback, not cached
- Sessions stored in DB, not JWT -- can be revoked by deleting the row
- No password storage, no password reset flow -- Google handles all credential management
