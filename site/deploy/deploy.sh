#!/bin/bash
set -euo pipefail

cd /opt/recipe-site/repo
git pull origin main
cd site
go build -o /opt/recipe-site/recipe-site .
sudo systemctl restart recipe-site
echo "Deployed successfully"
