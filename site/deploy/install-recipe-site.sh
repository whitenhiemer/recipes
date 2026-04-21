#!/usr/bin/env bash
#
# Recipe Site Install Script
# Runs inside the LXC container -- installs all dependencies, builds the app,
# configures nginx, systemd, and SQLite database.
#
set -eEuo pipefail

REPO_URL="${1:-https://github.com/whitenhiemer/recipes.git}"
APP_DIR="/opt/recipe-site"
REPO_DIR="${APP_DIR}/repo"
DATA_DIR="${APP_DIR}/data"
BIN_PATH="${APP_DIR}/recipe-site"
DB_PATH="${DATA_DIR}/recipes.db"

# -- Colors & logging ---------------------------------------------------------
BL='\033[36m'  YW='\033[33m'  RD='\033[91m'  GN='\033[92m'
WH='\033[97m'  CL='\033[0m'

info()  { echo -e "${BL}[INFO]${CL}  ${WH}$*${CL}"; }
warn()  { echo -e "${YW}[WARN]${CL}  ${WH}$*${CL}"; }
error() { echo -e "${RD}[ERROR]${CL} ${WH}$*${CL}"; exit 1; }
msg()   { echo -e "${GN}$*${CL}"; }

# -- System packages -----------------------------------------------------------
install_deps() {
    info "Updating package lists..."
    apt-get update -qq

    info "Installing dependencies..."
    apt-get install -y -qq \
        nginx \
        git \
        golang-go \
        sqlite3 \
        certbot \
        python3-certbot-nginx \
        curl \
        > /dev/null 2>&1

    info "Dependencies installed"
}

# -- Service user --------------------------------------------------------------
create_user() {
    if id recipe-site &>/dev/null; then
        info "Service user already exists"
        return
    fi

    info "Creating service user..."
    useradd -r -s /bin/false -m -d "$APP_DIR" recipe-site
}

# -- Directory structure -------------------------------------------------------
setup_dirs() {
    info "Setting up directories..."
    mkdir -p "$APP_DIR" "$DATA_DIR"
    chown -R recipe-site:recipe-site "$APP_DIR"
}

# -- Clone repo ----------------------------------------------------------------
clone_repo() {
    if [[ -d "$REPO_DIR/.git" ]]; then
        info "Repository already exists, pulling latest..."
        git -C "$REPO_DIR" pull origin main
    else
        info "Cloning repository..."
        git clone "$REPO_URL" "$REPO_DIR"
    fi
    chown -R recipe-site:recipe-site "$REPO_DIR"
}

# -- Build Go binary -----------------------------------------------------------
build_app() {
    info "Building Go application..."
    cd "$REPO_DIR/site"
    go build -o "$BIN_PATH" .
    chown recipe-site:recipe-site "$BIN_PATH"
    info "Binary built at $BIN_PATH"
}

# -- Webhook secret ------------------------------------------------------------
setup_webhook_secret() {
    local secret_file="${APP_DIR}/.webhook-secret"
    if [[ -f "$secret_file" ]]; then
        info "Webhook secret already exists"
        return
    fi

    info "Generating webhook secret..."
    openssl rand -hex 32 > "$secret_file"
    chmod 600 "$secret_file"
    chown recipe-site:recipe-site "$secret_file"
}

# -- SQLite database -----------------------------------------------------------
setup_database() {
    info "Initializing SQLite database..."
    touch "$DB_PATH"
    chown recipe-site:recipe-site "$DB_PATH"

    # Set up daily backup cron
    local backup_dir="${DATA_DIR}/backups"
    mkdir -p "$backup_dir"
    chown recipe-site:recipe-site "$backup_dir"

    cat > /etc/cron.d/recipe-site-backup <<CRON
# Daily SQLite backup at 3 AM
0 3 * * * recipe-site sqlite3 ${DB_PATH} ".backup '${backup_dir}/recipes-\$(date +\\%Y\\%m\\%d).db'" 2>/dev/null
# Keep 14 days of backups
5 3 * * * recipe-site find ${backup_dir} -name "recipes-*.db" -mtime +14 -delete 2>/dev/null
CRON
    chmod 644 /etc/cron.d/recipe-site-backup
    info "Database backup cron installed (daily, 14-day retention)"
}

# -- systemd service -----------------------------------------------------------
install_service() {
    info "Installing systemd service..."
    cat > /etc/systemd/system/recipe-site.service <<SERVICE
[Unit]
Description=Recipe Site - House Woodward Gourmand
After=network.target

[Service]
Type=simple
User=recipe-site
Group=recipe-site
WorkingDirectory=${REPO_DIR}/site
ExecStart=${BIN_PATH} \\
    -addr=127.0.0.1:8080 \\
    -recipes=${REPO_DIR} \\
    -db=${DB_PATH} \\
    -webhook-secret-file=${APP_DIR}/.webhook-secret
Restart=always
RestartSec=5

NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=${APP_DIR}
PrivateTmp=true
ProtectKernelTunables=true
ProtectControlGroups=true

Environment=GOMAXPROCS=1

[Install]
WantedBy=multi-user.target
SERVICE

    systemctl daemon-reload
    systemctl enable recipe-site
    systemctl start recipe-site
    info "Service installed and started"

    sleep 2
    if systemctl is-active --quiet recipe-site; then
        info "Service is running"
    else
        warn "Service may not have started correctly, check: journalctl -u recipe-site"
    fi
}

# -- nginx reverse proxy -------------------------------------------------------
configure_nginx() {
    info "Configuring nginx..."

    local server_ip
    server_ip=$(hostname -I | awk '{print $1}')

    cat > /etc/nginx/sites-available/recipe-site <<NGINX
server {
    listen 80;
    server_name ${server_ip} _;

    add_header X-Content-Type-Options "nosniff" always;
    add_header X-Frame-Options "DENY" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    location /static/ {
        alias ${REPO_DIR}/site/static/;
        expires 7d;
        add_header Cache-Control "public, immutable";
    }

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
NGINX

    ln -sf /etc/nginx/sites-available/recipe-site /etc/nginx/sites-enabled/
    rm -f /etc/nginx/sites-enabled/default

    nginx -t
    systemctl restart nginx
    info "nginx configured and running"
}

# -- Verify --------------------------------------------------------------------
verify_install() {
    info "Verifying installation..."

    local attempts=0
    while [[ $attempts -lt 5 ]]; do
        if curl -sf http://127.0.0.1:8080/ >/dev/null 2>&1; then
            info "Application is responding on port 8080"
            return 0
        fi
        ((attempts++))
        sleep 2
    done

    warn "Application not responding yet -- it may still be starting"
}

# -- MOTD ---------------------------------------------------------------------
set_motd() {
    local ip
    ip=$(hostname -I | awk '{print $1}')

    cat > /etc/motd <<MOTD

    House Woodward Gourmand - Recipe Site

    Access:    http://${ip}
    Service:   systemctl status recipe-site
    Logs:      journalctl -u recipe-site -f
    Database:  ${DB_PATH}
    Backups:   ${DATA_DIR}/backups/
    Deploy:    ${REPO_DIR}/site/deploy/deploy.sh

MOTD
}

# -- Main ----------------------------------------------------------------------
main() {
    msg ""
    msg "==========================================="
    msg "  Installing $APP"
    msg "==========================================="
    msg ""

    install_deps
    create_user
    setup_dirs
    clone_repo
    build_app
    setup_webhook_secret
    setup_database
    install_service
    configure_nginx
    verify_install
    set_motd

    msg ""
    msg "==========================================="
    msg "  Installation complete!"
    msg "==========================================="
    msg ""
}

main "$@"
