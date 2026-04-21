#!/bin/bash
set -euo pipefail

echo "=== Recipe Site LXC Setup ==="

apt-get update && apt-get install -y nginx git golang-go

useradd -r -s /bin/false -m -d /opt/recipe-site recipe-site || true

mkdir -p /opt/recipe-site
chown recipe-site:recipe-site /opt/recipe-site

echo "Clone the repo:"
echo "  git clone git@github.com:whitenhiemer/recipes.git /opt/recipe-site/repo"
echo ""
echo "Build:"
echo "  cd /opt/recipe-site/repo/site && go build -o /opt/recipe-site/recipe-site ."
echo ""
echo "Create webhook secret:"
echo "  openssl rand -hex 32 > /opt/recipe-site/.webhook-secret"
echo "  chmod 600 /opt/recipe-site/.webhook-secret"
echo "  chown recipe-site:recipe-site /opt/recipe-site/.webhook-secret"
echo ""
echo "Install service:"
echo "  cp /opt/recipe-site/repo/site/deploy/recipe-site.service /etc/systemd/system/"
echo "  systemctl daemon-reload"
echo "  systemctl enable --now recipe-site"
echo ""
echo "Install nginx:"
echo "  cp /opt/recipe-site/repo/site/deploy/nginx.conf /etc/nginx/sites-available/recipe-site"
echo "  ln -sf /etc/nginx/sites-available/recipe-site /etc/nginx/sites-enabled/"
echo "  rm -f /etc/nginx/sites-enabled/default"
echo "  systemctl restart nginx"
echo ""
echo "TLS (optional):"
echo "  apt-get install -y certbot python3-certbot-nginx"
echo "  certbot --nginx -d recipes.yourdomain.com"
