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

You can also add recipes directly from the website using the "+ New Recipe" button on the home page.

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

### 1. Create the LXC Container

From the Proxmox web UI or CLI, create a new container:

```bash
# From the Proxmox host
pct create <CTID> local:vztmpl/debian-12-standard_12.7-1_amd64.tar.zst \
  --hostname recipe-site \
  --memory 256 \
  --swap 256 \
  --rootfs local-lvm:2 \
  --cores 1 \
  --net0 name=eth0,bridge=vmbr0,ip=dhcp \
  --unprivileged 1 \
  --start 1
```

Adjust `CTID`, storage, and network bridge to match your environment. 256MB RAM and 2GB disk is sufficient.

### 2. Initial Setup Inside the Container

Attach to the container and install dependencies:

```bash
pct enter <CTID>

apt-get update && apt-get install -y nginx git golang-go

# Create a service user
useradd -r -s /bin/false -m -d /opt/recipe-site recipe-site
mkdir -p /opt/recipe-site
chown recipe-site:recipe-site /opt/recipe-site
```

### 3. Clone and Build

```bash
# Clone the repo (use HTTPS if SSH keys aren't set up in the container)
git clone https://github.com/whitenhiemer/recipes.git /opt/recipe-site/repo

# Build the Go binary
cd /opt/recipe-site/repo/site
go build -o /opt/recipe-site/recipe-site .

# Set ownership
chown recipe-site:recipe-site /opt/recipe-site/recipe-site
```

### 4. Create the Webhook Secret

This is used to verify GitHub webhook payloads:

```bash
openssl rand -hex 32 > /opt/recipe-site/.webhook-secret
chmod 600 /opt/recipe-site/.webhook-secret
chown recipe-site:recipe-site /opt/recipe-site/.webhook-secret
```

### 5. Install the systemd Service

```bash
cp /opt/recipe-site/repo/site/deploy/recipe-site.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable --now recipe-site

# Verify it's running
systemctl status recipe-site
curl -s http://127.0.0.1:8080/ | head -5
```

The service runs the Go binary on `127.0.0.1:8080` with security hardening (NoNewPrivileges, ProtectSystem=strict, read-only filesystem except the repo directory).

### 6. Configure nginx Reverse Proxy

```bash
# Edit the server_name in the config first
vim /opt/recipe-site/repo/site/deploy/nginx.conf
# Change "recipes.yourdomain.com" to your actual domain or the container's IP

cp /opt/recipe-site/repo/site/deploy/nginx.conf /etc/nginx/sites-available/recipe-site
ln -sf /etc/nginx/sites-available/recipe-site /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default
systemctl restart nginx
```

The site should now be accessible at `http://<container-ip>`.

### 7. TLS with Let's Encrypt (Optional)

If the container is publicly accessible with a domain name:

```bash
apt-get install -y certbot python3-certbot-nginx
certbot --nginx -d recipes.yourdomain.com
```

### 8. Configure GitHub Webhook for Auto-Reload

1. Go to the repo settings on GitHub: Settings > Webhooks > Add webhook
2. Set the payload URL to `https://recipes.yourdomain.com/webhook`
3. Set the content type to `application/json`
4. Set the secret to the contents of `/opt/recipe-site/.webhook-secret`
5. Select "Just the push event"

When you push new recipes, the webhook triggers a `git pull` and index reload -- no restart needed.

### 9. Manual Deploys (Code Changes)

For changes to the Go code, templates, or static assets, run the deploy script:

```bash
/opt/recipe-site/repo/site/deploy/deploy.sh
```

This pulls the latest code, rebuilds the binary, and restarts the service.

### Container Summary

| Setting | Value |
|---------|-------|
| OS | Debian 12 |
| RAM | 256 MB |
| Disk | 2 GB |
| CPU | 1 core |
| Stack | nginx -> Go binary on :8080 |
| Repo | `/opt/recipe-site/repo` |
| Binary | `/opt/recipe-site/recipe-site` |
| Service | `recipe-site.service` |
| Webhook secret | `/opt/recipe-site/.webhook-secret` |

## Project Plan

See [PLAN.md](PLAN.md) for architecture details and project roadmap.
