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
├── lunch/              # Midday meals and sandwiches
├── sauces-marinades/   # Sauces, dressings, rubs, brines, and marinades
├── sides/              # Side dishes and accompaniments
├── snacks/             # Appetizers and quick bites
├── soups-stews/        # Soups, stews, and chilis
├── TEMPLATE.md         # Recipe template
└── site/               # Go web application
    ├── main.go
    ├── Makefile
    ├── internal/       # Go packages (config, recipe parser, store, handlers)
    ├── templates/      # HTML templates (layout, home, recipe, etc.)
    ├── static/         # CSS, JS, images
    └── deploy/         # Proxmox deploy scripts, systemd, nginx
```

## Storage

Recipes are stored in a **SQLite database** (WAL mode). On startup, the app imports any markdown recipe files found in the repo into the database. New recipes created through the web UI are stored directly in the database.

**Backups:** A daily cron job runs `sqlite3 .backup` at 3 AM with 14-day retention. The database file lives at `/opt/recipe-site/data/recipes.db` on the deployed container.

**Migration path:** The schema is PostgreSQL-compatible. When the time comes, swap `modernc.org/sqlite` for `lib/pq` and update the connection string.

## Adding a Recipe

1. Copy [TEMPLATE.md](TEMPLATE.md) into the appropriate category folder
2. Rename it to match the recipe (e.g., `smoked-brisket.md`)
3. Fill in the sections following the template format
4. Commit and push -- the webhook will auto-reload the site

You can also add recipes directly from the website using the "+ New Recipe" button on the home page.

## Features

- **Search & browse** -- live search, category tabs, tag filtering
- **Meal planning** -- dropdown selects for 4 slots per meal, split into early/late week sections
- **Shopping list** -- auto-generated from meal plan, grouped by grocery store department with price estimates
- **Two-trip shopping** -- shelf-stable items in trip 1, perishables split by early/late week to reduce spoilage
- **Buy units** -- standard purchase sizes shown per ingredient (e.g., "1 lb / 4 sticks" for butter)
- **Pantry staples** -- common items (salt, oil, flour, spices) flagged and excluded from estimated total
- **Pantry & fridge inventory** -- track what you have on hand with prepopulated suggestion chips (emoji icons, grouped by location); on-hand items auto-detected on the shopping list via fuzzy matching
- **Print-friendly** -- branded single-page shopping list (side-by-side trips, compact layout with logo)
- **Recipe creation** -- step-by-step wizard, markdown file upload, or import from URL
- **URL import** -- paste a recipe URL to extract data via JSON-LD; source attribution on recipe detail pages
- **Inline recipe search** -- search for recipes online with top 3 suggestions and one-click fetch to import
- **Recipe scaling** -- default 5 servings with +/- control; ingredient quantities scale with fraction display
- **Kitchen mode** -- screen wake lock keeps your phone awake while cooking
- **Mobile-first** -- responsive design with touch-friendly controls

### Recipe Format

```markdown
---
tags: [tag1, tag2]
image: recipe-name.jpg
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

Tags and image are optional. Category is derived from the directory name. Subsection headers (`###`) can be used under Ingredients for multi-component recipes. Images go in `site/static/img/recipes/`.

## Running the Site Locally

```bash
cd site
make run
# Opens at http://localhost:8080
```

## Deploying to Proxmox LXC

### Automated Deploy (Recommended)

The deploy scripts follow the [community-scripts.org](https://community-scripts.org) convention. From your Proxmox host:

```bash
# Copy the deploy scripts to your Proxmox host, then run:
bash ct-recipe-site.sh
```

This launches a whiptail wizard that:
1. Prompts for container ID, resources, and network settings (or uses defaults)
2. Creates the LXC container with Debian 12
3. Installs all dependencies (nginx, Go, SQLite)
4. Clones the repo and builds the Go binary
5. Configures systemd, nginx reverse proxy, and database backups
6. Generates a webhook secret

The scripts are in `site/deploy/`:
- `ct-recipe-site.sh` -- runs on Proxmox host, creates the LXC
- `install-recipe-site.sh` -- runs inside the container, installs everything

### Manual Deploy

If you prefer to set things up manually:

#### 1. Create the LXC Container

```bash
pct create <CTID> local:vztmpl/debian-12-standard_12.7-1_amd64.tar.zst \
  --hostname recipe-site \
  --memory 512 \
  --swap 512 \
  --rootfs local-lvm:4 \
  --cores 1 \
  --net0 name=eth0,bridge=vmbr0,ip=dhcp \
  --unprivileged 1 \
  --features nesting=1 \
  --start 1
```

#### 2. Install Dependencies

```bash
pct enter <CTID>
apt-get update && apt-get install -y nginx git golang-go sqlite3
```

#### 3. Clone, Build, and Configure

```bash
useradd -r -s /bin/false -m -d /opt/recipe-site recipe-site
mkdir -p /opt/recipe-site/data
git clone https://github.com/whitenhiemer/recipes.git /opt/recipe-site/repo
cd /opt/recipe-site/repo/site && go build -o /opt/recipe-site/recipe-site .
chown -R recipe-site:recipe-site /opt/recipe-site
openssl rand -hex 32 > /opt/recipe-site/.webhook-secret
chmod 600 /opt/recipe-site/.webhook-secret
chown recipe-site:recipe-site /opt/recipe-site/.webhook-secret
```

#### 4. Install Service and nginx

```bash
cp /opt/recipe-site/repo/site/deploy/recipe-site.service /etc/systemd/system/
systemctl daemon-reload && systemctl enable --now recipe-site

# Edit server_name in nginx.conf, then:
cp /opt/recipe-site/repo/site/deploy/nginx.conf /etc/nginx/sites-available/recipe-site
ln -sf /etc/nginx/sites-available/recipe-site /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default && systemctl restart nginx
```

#### 5. TLS (Optional)

```bash
apt-get install -y certbot python3-certbot-nginx
certbot --nginx -d recipes.yourdomain.com
```

### GitHub Webhook for Auto-Reload

1. Go to repo Settings > Webhooks > Add webhook
2. Payload URL: `https://recipes.yourdomain.com/webhook`
3. Content type: `application/json`
4. Secret: contents of `/opt/recipe-site/.webhook-secret`
5. Events: "Just the push event"

On push, the webhook triggers `git pull`, reimports markdown into the DB, and rebuilds the index.

### Manual Code Deploys

```bash
/opt/recipe-site/repo/site/deploy/deploy.sh
```

### Container Summary

| Setting | Value |
|---------|-------|
| OS | Debian 12 |
| RAM | 512 MB |
| Disk | 4 GB |
| CPU | 1 core |
| Stack | nginx -> Go binary on :8080 -> SQLite |
| Repo | `/opt/recipe-site/repo` |
| Database | `/opt/recipe-site/data/recipes.db` |
| Backups | `/opt/recipe-site/data/backups/` (daily, 14-day retention) |
| Binary | `/opt/recipe-site/recipe-site` |
| Service | `recipe-site.service` |
| Webhook secret | `/opt/recipe-site/.webhook-secret` |

## Project Plan

See [PLAN.md](PLAN.md) for architecture details and project roadmap.
