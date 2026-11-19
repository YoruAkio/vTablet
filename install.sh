#!/usr/bin/env bash
#
# vTablet Linux Installer
# Installs vTablet server, udev rules for /dev/uinput, and optional systemd user service.
#
set -euo pipefail

REPO="YoruAkio/vTablet"
INSTALL_DIR="/usr/local/bin"
VARIANT="headless" # "headless" (recommended, zero-CGO) or "gui" (system tray)
ENABLE_SERVICE=false
SYSTEMD_DIR="${HOME}/.config/systemd/user"

log_info() {
    printf "[INFO] %s\n" "$1"
}

log_warn() {
    printf "[WARN] %s\n" "$1" >&2
}

log_err() {
    printf "[ERROR] %s\n" "$1" >&2
}

run_sudo() {
    if [ "$(id -u)" -eq 0 ]; then
        "$@"
    elif command -v sudo >/dev/null 2>&1; then
        sudo "$@"
    else
        log_err "Root privileges required for this step, but 'sudo' was not found."
        return 1
    fi
}

show_help() {
    cat <<EOF
vTablet Linux Installer

Usage:
  ./install.sh [options]

Options:
  --gui             Install the GUI version (with system tray indicator; requires GTK)
  --headless        Install the Headless version (default; zero external dependencies)
  --user            Install binary to ~/.local/bin instead of /usr/local/bin
  --dir <path>      Specify custom installation directory for binary
  --service         Enable and start systemd user service immediately
  --help, -h        Show this help message
EOF
}

# Parse CLI arguments
while [ $# -gt 0 ]; do
    case "$1" in
        --gui)
            VARIANT="gui"
            shift
            ;;
        --headless)
            VARIANT="headless"
            shift
            ;;
        --user)
            INSTALL_DIR="${HOME}/.local/bin"
            shift
            ;;
        --dir)
            if [ -n "${2:-}" ]; then
                INSTALL_DIR="$2"
                shift 2
            else
                log_err "--dir requires a path argument"
                exit 1
            fi
            ;;
        --service)
            ENABLE_SERVICE=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            log_warn "Unknown option: $1"
            shift
            ;;
    esac
done

ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64)
        PKG_ARCH="amd64"
        ;;
    *)
        log_err "Unsupported architecture: $ARCH (currently only x86_64 / amd64 binaries are provided in releases)."
        exit 1
        ;;
esac

log_info "Target architecture: $PKG_ARCH"
log_info "Target installation directory: $INSTALL_DIR"
log_info "Selected binary variant: $VARIANT"

# Create a temporary workspace for downloads
TMP_DIR="$(mktemp -d -t vtablet-install.XXXXXX)"
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

# 1. Obtain binary
TARGET_BIN="$TMP_DIR/vtablet"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" 2>/dev/null && pwd || echo "")"

# Check if running within a cloned repository with existing local build
if [ -n "$SCRIPT_DIR" ] && [ -f "$SCRIPT_DIR/server/vtablet" ]; then
    log_info "Found existing local build in repository at $SCRIPT_DIR/server/vtablet."
    cp "$SCRIPT_DIR/server/vtablet" "$TARGET_BIN"
else
    log_info "Downloading latest $VARIANT release from GitHub ($REPO)..."
    if [ "$VARIANT" = "gui" ]; then
        TAR_NAME="vtablet-linux-${PKG_ARCH}.tar.gz"
    else
        TAR_NAME="vtablet-linux-headless-${PKG_ARCH}.tar.gz"
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${TAR_NAME}"
    TAR_DEST="$TMP_DIR/$TAR_NAME"

    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$DOWNLOAD_URL" -o "$TAR_DEST"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$DOWNLOAD_URL" -O "$TAR_DEST"
    else
        log_err "Neither curl nor wget was found. Please install one to proceed."
        exit 1
    fi

    tar -xzf "$TAR_DEST" -C "$TMP_DIR"
    if [ ! -f "$TARGET_BIN" ]; then
        log_err "Extracted archive did not contain expected 'vtablet' executable."
        exit 1
    fi
fi

# 2. Install binary
chmod +x "$TARGET_BIN"
if [ "$INSTALL_DIR" = "/usr/local/bin" ] && [ "$(id -u)" -ne 0 ]; then
    log_info "Installing vtablet to $INSTALL_DIR (requires sudo)..."
    run_sudo install -d "$INSTALL_DIR"
    run_sudo install -m 755 "$TARGET_BIN" "$INSTALL_DIR/vtablet"
else
    mkdir -p "$INSTALL_DIR"
    install -m 755 "$TARGET_BIN" "$INSTALL_DIR/vtablet"
fi
INSTALLED_BIN="$INSTALL_DIR/vtablet"
log_info "Successfully installed binary to $INSTALLED_BIN"

# 3. Configure udev rules for /dev/uinput
UDEV_RULE_FILE="/etc/udev/rules.d/99-vtablet.rules"
log_info "Configuring udev rules for /dev/uinput..."

TMP_UDEV="$TMP_DIR/99-vtablet.rules"
cat <<'EOF' > "$TMP_UDEV"
# udev rule to allow non-root users in 'input' group or active session to access /dev/uinput
KERNEL=="uinput", MODE="0660", GROUP="input", OPTIONS+="static_node=uinput", TAG+="uaccess"
EOF

run_sudo cp "$TMP_UDEV" "$UDEV_RULE_FILE"
run_sudo chmod 644 "$UDEV_RULE_FILE"

if command -v udevadm >/dev/null 2>&1; then
    run_sudo udevadm control --reload-rules || true
    run_sudo udevadm trigger || true
fi
log_info "udev rules installed at $UDEV_RULE_FILE"

# 4. User group membership
CURRENT_USER="${SUDO_USER:-$USER}"
if [ -n "$CURRENT_USER" ] && [ "$CURRENT_USER" != "root" ]; then
    if id -nG "$CURRENT_USER" | grep -qw "input"; then
        log_info "User '$CURRENT_USER' is already a member of the 'input' group."
    else
        log_info "Adding user '$CURRENT_USER' to the 'input' group..."
        run_sudo usermod -aG input "$CURRENT_USER" || log_warn "Could not add $CURRENT_USER to input group automatically."
        log_warn "You must log out and log back in (or reboot) for group changes to take effect."
    fi
fi

# 5. Optional systemd user service
if command -v systemctl >/dev/null 2>&1; then
    mkdir -p "$SYSTEMD_DIR"
    SERVICE_FILE="$SYSTEMD_DIR/vtablet.service"

    cat <<EOF > "$SERVICE_FILE"
[Unit]
Description=vTablet Virtual Digitizer Server
After=graphical-session.target

[Service]
ExecStart=${INSTALLED_BIN} -headless
Restart=on-failure
RestartSec=3

[Install]
WantedBy=default.target
EOF

    systemctl --user daemon-reload || true
    log_info "systemd user service installed at $SERVICE_FILE"

    if [ "$ENABLE_SERVICE" = true ]; then
        log_info "Enabling and starting vtablet systemd user service..."
        systemctl --user enable --now vtablet.service
        log_info "Service is running! Check status with: systemctl --user status vtablet.service"
    else
        log_info "To start vTablet as a background service on login:"
        log_info "  systemctl --user enable --now vtablet.service"
    fi
fi

# Check PATH
case ":$PATH:" in
    *:"$INSTALL_DIR":*) ;;
    *)
        log_warn "$INSTALL_DIR is not in your PATH. You may want to add it to your shell configuration."
        ;;
esac

log_info "Installation complete!"
log_info "Run manually with:"
log_info "  $INSTALLED_BIN -headless"
