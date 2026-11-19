# vTablet

> Use your tablet or smartphone as a low-latency virtual digitizer / graphics tablet for *osu!* and digital art!

[![Build Golang Server](https://github.com/Teages/vTablet/actions/workflows/build-server.yml/badge.svg)](https://github.com/Teages/vTablet/actions/workflows/build-server.yml)

> [!WARNING]
> 3.0 is currently in preview version, and there may be incompatible changes in the transmission protocol. Please be sure to update the server and client at the same time.

---

## Features

- **Cross-Platform Host Server:**
  - **Linux (Wayland & Hyprland first-class):** Emulates a native graphics tablet pen tool via Linux kernel `/dev/uinput` with sub-millisecond latency. Works natively with `libinput`, Hyprland, Sway, KDE, GNOME, and *osu!lazer* / *osu!* under Wine.
  - **Windows:** Emulates native Windows pointer devices via Win32 Synthetic Pointer APIs.
- **Cross-Platform Client (Flutter):**
  - Android APK, iOS, and Web PWA.
  - Pressure sensitivity (up to 8192 levels) and tilt support.
  - Automatic screen mapping and aspect-ratio preservation.
- **Automatic USB Tethering (ADB):**
  - Seamless plug-and-play reverse port forwarding over USB via ADB.
- **Headless Daemon Support:**
  - Can run in headless mode as a background `systemd` user service without requiring a GUI or X11 tray.

---

## Linux Setup (Wayland / Hyprland)

### 1. Permissions (`/dev/uinput`)

To allow the server to create virtual tablet input devices without running as root:

```bash
# Copy udev rule
sudo cp server/99-vtablet.rules /etc/udev/rules.d/

# Ensure your user is in the 'input' group
sudo usermod -aG input $USER

# Reload udev rules
sudo udevadm control --reload-rules && sudo udevadm trigger
```

### 2. Building & Running the Server

```bash
cd server
go build -o vtablet .

# Run in headless mode (recommended for Wayland / Hyprland)
./vtablet -headless

# Or run with system tray (if supported by your desktop environment)
./vtablet
```

### 3. Hyprland Configuration

In `hyprland.conf`, you can optionally bind the virtual tablet to a specific output / monitor:

```ini
# ~/.config/hypr/hyprland.conf
device {
    name = vtablet-digitizer
    output = eDP-1   # replace with your target monitor (e.g. DP-1, HDMI-A-1)
}
```

### 4. Running as a `systemd` User Service

```bash
mkdir -p ~/.config/systemd/user
cp server/vtablet.service ~/.config/systemd/user/
systemctl --user daemon-reload
systemctl --user enable --now vtablet.service
```

---

## Client Setup

- **Android:** Install the APK built from `client/` or open the Web app at `https://vtablet.teages.xyz`.
- **USB Connection (Fastest & Lowest Latency):**
  1. Enable USB Debugging on your Android tablet.
  2. Connect tablet via USB. The server will automatically run `adb reverse` and launch the client.
- **Wi-Fi Connection:**
  1. Ensure both devices are on the same local network.
  2. In the client, enter your host IP and port `23052` (e.g., `192.168.1.100:23052`).

---

## Protocol

Input frames are transmitted over WebSocket in compact 8-byte binary packets:

```
[2 bytes Type] [2 bytes Code] [4 bytes Value]
```

- `0x0003` (`EV_ABS`): `0x00` (`ABS_X`), `0x01` (`ABS_Y`), `0x18` (`ABS_PRESSURE`), `0x1A` (`ABS_TILT_X`), `0x1B` (`ABS_TILT_Y`)
- `0x0000` (`EV_SYN`): `0x00` (`SYN_REPORT`), `0xFFFF` (`SYN_PING`)
