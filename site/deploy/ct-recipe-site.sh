#!/usr/bin/env bash
#
# Recipe Site LXC Container Creator
# Runs on the Proxmox host to create and configure an LXC container
# Style: community-scripts.org convention
#
set -eEuo pipefail

# -- Colors & logging ---------------------------------------------------------
BL='\033[36m'  YW='\033[33m'  RD='\033[91m'  GN='\033[92m'
WH='\033[97m'  CL='\033[0m'

info()  { echo -e "${BL}[INFO]${CL}  ${WH}$*${CL}"; }
warn()  { echo -e "${YW}[WARN]${CL}  ${WH}$*${CL}"; }
error() { echo -e "${RD}[ERROR]${CL} ${WH}$*${CL}"; }
msg()   { echo -e "${GN}$*${CL}"; }

die() {
    local exit_code=$?
    error "Script failed at line $BASH_LINENO with exit code $exit_code"
    cleanup_on_failure
    exit "$exit_code"
}
trap die ERR

# -- Defaults ------------------------------------------------------------------
APP="Recipe Site"
APP_TAG="recipe-site"
REPO_URL="https://github.com/whitenhiemer/recipes.git"

CT_ID=""
CT_HOSTNAME="recipe-site"
CT_DISK="4"
CT_RAM="512"
CT_CORES="1"
CT_BRIDGE="vmbr0"
CT_STORAGE="local-lvm"
CT_TEMPLATE_STORAGE="local"
CT_OS_TYPE="debian"
CT_OS_VERSION="12"
CT_UNPRIVILEGED="1"

CREATED_CT=""

# -- Cleanup -------------------------------------------------------------------
cleanup_on_failure() {
    if [[ -n "$CREATED_CT" ]]; then
        warn "Cleaning up container $CREATED_CT..."
        pct stop "$CREATED_CT" 2>/dev/null || true
        pct destroy "$CREATED_CT" 2>/dev/null || true
    fi
}

# -- Validation helpers --------------------------------------------------------
validate_ct_id() {
    local id=$1
    if ! [[ "$id" =~ ^[0-9]+$ ]]; then
        return 1
    fi
    if [[ -f "/etc/pve/lxc/${id}.conf" ]] || [[ -f "/etc/pve/qemu-server/${id}.conf" ]]; then
        return 1
    fi
    return 0
}

next_ct_id() {
    local id
    id=$(pvesh get /cluster/nextid 2>/dev/null) || id=100
    echo "$id"
}

# -- Storage helpers -----------------------------------------------------------
get_storage_list() {
    pvesm status -content rootdir 2>/dev/null | awk 'NR>1 {print $1}' | sort
}

get_template_storage_list() {
    pvesm status -content vztmpl 2>/dev/null | awk 'NR>1 {print $1}' | sort
}

get_template_path() {
    local storage=$1
    local tmpl
    tmpl=$(pveam list "$storage" 2>/dev/null | grep "debian-12" | sort -V | tail -1 | awk '{print $1}')
    if [[ -z "$tmpl" ]]; then
        info "Downloading Debian 12 template..."
        pveam update >/dev/null 2>&1
        local available
        available=$(pveam available --section system 2>/dev/null | grep "debian-12" | sort -V | tail -1 | awk '{print $2}')
        if [[ -z "$available" ]]; then
            error "No Debian 12 template found"
            exit 1
        fi
        pveam download "$storage" "$available"
        tmpl="${storage}:vztmpl/${available}"
    fi
    echo "$tmpl"
}

# -- Whiptail dialogs ----------------------------------------------------------
check_whiptail() {
    if ! command -v whiptail &>/dev/null; then
        error "whiptail is required but not installed"
        exit 1
    fi
}

dialog_yesno() {
    whiptail --backtitle "Proxmox VE Helper Scripts" \
        --title "$APP LXC" \
        --yesno "$1" 10 60
}

dialog_input() {
    local title=$1 prompt=$2 default=$3
    whiptail --backtitle "Proxmox VE Helper Scripts" \
        --title "$title" \
        --inputbox "$prompt" 10 60 "$default" 3>&1 1>&2 2>&3
}

dialog_menu() {
    local title=$1 prompt=$2
    shift 2
    whiptail --backtitle "Proxmox VE Helper Scripts" \
        --title "$title" \
        --menu "$prompt" 16 60 6 "$@" 3>&1 1>&2 2>&3
}

# -- Main flow -----------------------------------------------------------------
header() {
    clear
    cat <<'BANNER'
    ____            _                _____ _ __
   / __ \___  _____(_)___  ___      / ___/(_) /____
  / /_/ / _ \/ ___/ / __ \/ _ \     \__ \/ / __/ _ \
 / _, _/  __/ /__/ / /_/ /  __/    ___/ / / /_/  __/
/_/ |_|\___/\___/_/ .___/\___/    /____/_/\__/\___/
                 /_/
    House Woodward Gourmand
BANNER
    echo ""
}

confirm_creation() {
    if ! dialog_yesno "This will create a new LXC container for $APP.\n\nProceed?"; then
        info "Cancelled."
        exit 0
    fi
}

configure_settings() {
    local use_defaults
    if dialog_yesno "Use default settings?\n\nCPU: ${CT_CORES} core  RAM: ${CT_RAM}MB  Disk: ${CT_DISK}GB\nOS: Debian 12  Bridge: ${CT_BRIDGE}\n\nSelect 'No' for advanced configuration."; then
        use_defaults="yes"
    else
        use_defaults="no"
    fi

    # Container ID
    local default_id
    default_id=$(next_ct_id)
    CT_ID=$(dialog_input "Container ID" "Enter container ID:" "$default_id") || exit 0
    if ! validate_ct_id "$CT_ID"; then
        error "Container ID $CT_ID is invalid or already in use"
        exit 1
    fi

    if [[ "$use_defaults" == "no" ]]; then
        CT_HOSTNAME=$(dialog_input "Hostname" "Enter hostname:" "$CT_HOSTNAME") || exit 0
        CT_CORES=$(dialog_input "CPU Cores" "Enter number of CPU cores:" "$CT_CORES") || exit 0
        CT_RAM=$(dialog_input "RAM (MB)" "Enter RAM in MB:" "$CT_RAM") || exit 0
        CT_DISK=$(dialog_input "Disk (GB)" "Enter disk size in GB:" "$CT_DISK") || exit 0

        # Storage selection
        local storages
        storages=$(get_storage_list)
        if [[ -n "$storages" ]]; then
            local menu_args=()
            while IFS= read -r s; do
                menu_args+=("$s" "")
            done <<< "$storages"
            CT_STORAGE=$(dialog_menu "Storage" "Select storage for container rootfs:" "${menu_args[@]}") || exit 0
        fi

        CT_BRIDGE=$(dialog_input "Network Bridge" "Enter network bridge:" "$CT_BRIDGE") || exit 0
    else
        CT_ID="$default_id"
        if ! validate_ct_id "$CT_ID"; then
            CT_ID=$(next_ct_id)
        fi
    fi
}

create_container() {
    info "Resolving Debian 12 template..."
    local tmpl
    tmpl=$(get_template_path "$CT_TEMPLATE_STORAGE")

    info "Creating LXC container $CT_ID ($CT_HOSTNAME)..."
    pct create "$CT_ID" "$tmpl" \
        --hostname "$CT_HOSTNAME" \
        --memory "$CT_RAM" \
        --swap "$CT_RAM" \
        --rootfs "${CT_STORAGE}:${CT_DISK}" \
        --cores "$CT_CORES" \
        --net0 "name=eth0,bridge=${CT_BRIDGE},ip=dhcp" \
        --unprivileged "$CT_UNPRIVILEGED" \
        --features nesting=1 \
        --onboot 1 \
        --tags "$APP_TAG" \
        --start 0

    CREATED_CT="$CT_ID"
    info "Container $CT_ID created"
}

start_container() {
    info "Starting container $CT_ID..."
    pct start "$CT_ID"
    sleep 3
}

wait_for_network() {
    info "Waiting for network..."
    local attempts=0
    local max_attempts=10
    while [[ $attempts -lt $max_attempts ]]; do
        if pct exec "$CT_ID" -- ping -c 1 -W 2 deb.debian.org &>/dev/null; then
            info "Network is up"
            return 0
        fi
        ((attempts++))
        sleep 3
    done
    error "Network not available after ${max_attempts} attempts"
    exit 1
}

run_install_script() {
    info "Copying install script into container..."
    local script_dir
    script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local install_script="${script_dir}/install-recipe-site.sh"

    if [[ ! -f "$install_script" ]]; then
        error "Install script not found at $install_script"
        exit 1
    fi

    pct push "$CT_ID" "$install_script" /tmp/install-recipe-site.sh
    pct exec "$CT_ID" -- chmod +x /tmp/install-recipe-site.sh
    pct exec "$CT_ID" -- bash /tmp/install-recipe-site.sh "$REPO_URL"
}

get_container_ip() {
    local ip=""
    local attempts=0
    while [[ $attempts -lt 5 ]]; do
        ip=$(pct exec "$CT_ID" -- hostname -I 2>/dev/null | awk '{print $1}')
        if [[ -n "$ip" ]]; then
            echo "$ip"
            return 0
        fi
        ((attempts++))
        sleep 3
    done
    echo "unknown"
}

show_completion() {
    local ip
    ip=$(get_container_ip)

    msg ""
    msg "================================================================"
    msg "  $APP installation complete!"
    msg "================================================================"
    msg ""
    msg "  Container ID:  $CT_ID"
    msg "  Hostname:      $CT_HOSTNAME"
    msg "  IP Address:    $ip"
    msg ""
    msg "  Access:        http://${ip}"
    msg ""
    msg "  Database:      /opt/recipe-site/data/recipes.db"
    msg "  Service:       systemctl status recipe-site"
    msg "  Logs:          journalctl -u recipe-site -f"
    msg ""
    msg "  To configure GitHub webhook for auto-reload:"
    msg "    1. Get the secret: cat /opt/recipe-site/.webhook-secret"
    msg "    2. Add webhook in GitHub repo settings"
    msg "    3. URL: http://${ip}/webhook"
    msg ""
    msg "  To set up TLS (optional):"
    msg "    pct exec $CT_ID -- certbot --nginx -d your.domain.com"
    msg ""
    msg "================================================================"
    CREATED_CT=""
}

# -- Entry point ---------------------------------------------------------------
main() {
    header
    check_whiptail
    confirm_creation
    configure_settings
    create_container
    start_container
    wait_for_network
    run_install_script
    show_completion
}

main "$@"
