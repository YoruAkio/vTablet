<div align="center">

  <img src="client/web/icons/Icon-192.png" alt="vTablet Logo" width="120" height="120" />

  # vTablet

  <p><strong>Transform your smartphone or tablet into a low-latency virtual graphics tablet and digitizer for <em>osu!</em> and digital art.</strong></p>

  <p>
    <a href="https://github.com/YoruAkio/vTablet/releases/latest">
      <img src="https://img.shields.io/github/v/release/YoruAkio/vTablet?color=3b82f6&logo=github&label=Release" alt="Latest Release" />
    </a>
    <a href="https://github.com/YoruAkio/vTablet/actions/workflows/build-server.yml">
      <img src="https://github.com/YoruAkio/vTablet/actions/workflows/build-server.yml/badge.svg" alt="Build Golang Server" />
    </a>
    <a href="https://github.com/YoruAkio/vTablet/actions/workflows/build-client.yml">
      <img src="https://github.com/YoruAkio/vTablet/actions/workflows/build-client.yml/badge.svg" alt="Build Flutter Client" />
    </a>
    <a href="https://github.com/YoruAkio/vTablet/actions/workflows/release.yml">
      <img src="https://github.com/YoruAkio/vTablet/actions/workflows/release.yml/badge.svg" alt="Release All Targets" />
    </a>
    <a href="https://github.com/YoruAkio/vTablet/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/YoruAkio/vTablet?color=10b981" alt="License" />
    </a>
  </p>

  <p>
    <a href="#quick-downloads">Quick Downloads</a> •
    <a href="#installation">Installation</a> •
    <a href="#features">Features</a> •
    <a href="#linux-setup-wayland--hyprland">Linux Setup</a> •
    <a href="#windows-setup">Windows Setup</a> •
    <a href="#client-setup-android--web">Client Setup</a> •
    <a href="#recommended-settings-for-osu">osu! Guide</a>
  </p>

</div>

---

## Quick Downloads

Pre-compiled binaries and installable APK packages are available from the [Latest Release](https://github.com/YoruAkio/vTablet/releases/latest):

| Target | File | Description |
| :--- | :--- | :--- |
| **Android Client** | [`vTablet-client-android.apk`](https://github.com/YoruAkio/vTablet/releases/latest) | Installable APK for Android phones and tablets |
| **Android Bundle** | [`vTablet-client-android.aab`](https://github.com/YoruAkio/vTablet/releases/latest) | Android App Bundle |
| **Linux Server (GUI)** | [`vtablet-linux-amd64.tar.gz`](https://github.com/YoruAkio/vTablet/releases/latest) | Linux x64 with System Tray (includes udev rules and systemd service) |
| **Linux Server (Headless)** | [`vtablet-linux-headless-amd64.tar.gz`](https://github.com/YoruAkio/vTablet/releases/latest) | Pure Go, Zero-CGO headless binary (ideal for background services) |
| **Windows Server** | [`vtablet-windows-amd64.zip`](https://github.com/YoruAkio/vTablet/releases/latest) | Windows 64-bit host server |

---

## Installation

### Linux (Automated Bash Installer)

Run the automated installer to download the binary, configure `/dev/uinput` udev rules, set user group permissions, and install the systemd user service:

```bash
curl -fsSL https://raw.githubusercontent.com/YoruAkio/vTablet/main/install.sh | bash
```

Or from a local clone of this repository:

```bash
chmod +x install.sh
./install.sh
```

Available installer options:
- `--headless`: Install the headless server with zero external shared library dependencies (default).
- `--gui`: Install the GUI version with system tray integration (requires GTK and libayatana-appindicator).
- `--user`: Install binary to `~/.local/bin` instead of `/usr/local/bin`.
- `--service`: Enable and start the systemd user service immediately.

---

### Windows (Automated PowerShell Installer)

Open PowerShell and run the installer:

```powershell
irm https://raw.githubusercontent.com/YoruAkio/vTablet/main/install.ps1 | iex
```

Or from a local clone of this repository:

```powershell
powershell -ExecutionPolicy Bypass -File .\install.ps1
```

The script extracts the latest 64-bit release into `%LOCALAPPDATA%\vTablet`, adds the directory to your user `PATH`, and creates a Start Menu shortcut.

---

## Features

- **Linux Support (Wayland & Hyprland):**
  - Direct kernel-level absolute pointer emulation via `/dev/uinput`.
  - Native compatibility with Mango WM, Hyprland, Sway, KDE Plasma, and GNOME.
  - Sub-millisecond input latency with atomic event batching.
- **Windows Host Support:**
  - Emulates native pointer devices via Windows Synthetic Pointer APIs.
- **Modern Client (Flutter 3):**
  - Multilingual support: English (Default), Bahasa Indonesia (ID), and Chinese (ZH).
  - In-app language switcher under the Appearance tab.
  - Up to 8192 levels of pressure sensitivity and tilt support.
  - Capacitive phone touch fallback pressure for non-stylus devices.
  - Active area aspect-ratio scaling and offset customization.
- **Automatic USB Tethering (ADB):**
  - Automatic reverse port forwarding (`adb reverse`) for zero-configuration, plug-and-play USB connection.
- **Headless Daemon:**
  - Lightweight headless mode (`-headless`) for running as a background systemd user service without requiring a GUI tray.

---

## Linux Setup (Wayland / Hyprland)

### 1. Manual `/dev/uinput` Permissions (Alternative to Installer)

If you prefer configuring udev rules manually without running `install.sh`:

```bash
# 1. Install udev rule
sudo cp server/99-vtablet.rules /etc/udev/rules.d/

# 2. Add your current user to the 'input' group
sudo usermod -aG input $USER

# 3. Reload udev rules
sudo udevadm control --reload-rules && sudo udevadm trigger
```

> [!NOTE]
> Log out and log back in (or reboot) for the group membership change to take effect.

### 2. Running the Server

```bash
# Headless mode (recommended for Wayland / tiling window managers)
vtablet -headless

# Headless mode with verbose packet logging
vtablet -headless -verbose

# GUI mode with system tray (requires libayatana-appindicator / GTK)
vtablet
```

### 3. Hyprland Configuration (Optional)

In your `hyprland.conf`, you can optionally bind the virtual tablet to a specific display output:

```ini
# ~/.config/hypr/hyprland.conf
device {
    name = vtablet-pointer
    output = eDP-1   # Replace with your monitor name (e.g. DP-1, HDMI-A-1)
}
```

### 4. Running as a `systemd` User Service

To have the server start automatically when your graphical session starts:

```bash
mkdir -p ~/.config/systemd/user
cp server/vtablet.service ~/.config/systemd/user/
systemctl --user daemon-reload
systemctl --user enable --now vtablet.service
```

---

## Windows Setup

1. Download [`vtablet-windows-amd64.zip`](https://github.com/YoruAkio/vTablet/releases/latest) from the Releases page (or use `install.ps1`).
2. Extract the archive and launch `vtablet.exe`.
3. The server will start and minimize to the system tray. To run headless, launch it from PowerShell or Command Prompt with `-headless`:
   ```powershell
   vtablet.exe -headless
   ```

---

## Client Setup (Android & Web)

### USB Connection (Fastest & Lowest Latency)

1. Enable **USB Debugging** in Android Developer Options on your mobile device.
2. Connect your device to your PC with a USB cable.
3. Start the `vtablet` server. The server will automatically detect the device, configure `adb reverse tcp:23052 tcp:23052`, and launch the client web app or native application.
4. Open the installed **vTablet Android App** and tap **Connect**.

### Wi-Fi Connection

1. Ensure both your PC and phone/tablet are connected to the same local network.
2. Open the vTablet app and enter your PC's local IP address and port `23052` (e.g., `192.168.1.100:23052`).
3. Tap **Connect**.

---

## Recommended Settings for osu!

When playing *osu!* on Linux (via Wine / Proton / Lutris):

> [!IMPORTANT]
> - **Raw Input:** Set **Raw Input to OFF** in *osu!* settings. Under Wine, Raw Input listens exclusively for relative mouse packets and will freeze or ignore absolute tablet/pointer warping.
> - **Mouse Sensitivity:** Set to **1.0x** for true 1:1 tablet mapping.

---

## Protocol

Packets are transmitted over WebSocket in compact 8-byte big-endian binary frames:

```text
[ 2 bytes Type ] [ 2 bytes Code ] [ 4 bytes Value ]
```

- `0x0003` (`EV_ABS`):
  - `0x00` (`ABS_X`): 0–32767
  - `0x01` (`ABS_Y`): 0–32767
  - `0x18` (`ABS_PRESSURE`): 0–8192
  - `0x1A` (`ABS_TILT_X`): -90–90
  - `0x1B` (`ABS_TILT_Y`): -90–90
- `0x0000` (`EV_SYN`):
  - `0x00` (`SYN_REPORT`): Dispatch frame
  - `0xFFFF` (`SYN_PING`): Latency measurement

---

## License

This project is open-source software licensed under the [MIT License](LICENSE).
