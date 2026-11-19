# vTablet Windows Installer
# Downloads and installs the vTablet host server for 64-bit Windows.
[CmdletBinding()]
param(
    [string]$InstallDir = (Join-Path $env:LOCALAPPDATA "vTablet"),
    [switch]$NoPathUpdate,
    [switch]$NoShortcut,
    [switch]$Force
)

$ErrorActionPreference = "Stop"

function Write-Info([string]$Message) {
    Write-Host "[INFO] $Message" -ForegroundColor Cyan
}

function Write-Warn([string]$Message) {
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Err([string]$Message) {
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# 1. Architecture check
if (-not [Environment]::Is64BitOperatingSystem) {
    Write-Err "vTablet requires a 64-bit Windows operating system."
    exit 1
}

$Repo = "YoruAkio/vTablet"
$ZipName = "vtablet-windows-amd64.zip"
$DownloadUrl = "https://github.com/$Repo/releases/latest/download/$ZipName"

Write-Info "Target installation directory: $InstallDir"

# 2. Check if already running
$runningProcesses = Get-Process -Name "vtablet*" -ErrorAction SilentlyContinue
if ($runningProcesses) {
    if ($Force) {
        Write-Info "Stopping existing vTablet processes..."
        $runningProcesses | Stop-Process -Force
    } else {
        Write-Warn "vTablet appears to be running. Please close it or run with -Force to terminate automatically."
        exit 1
    }
}

# 3. Download release zip
$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("vtablet-install-" + [System.Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

try {
    $zipPath = Join-Path $tempDir $ZipName
    Write-Info "Downloading latest release from $DownloadUrl..."

    # Use TLS 1.2+ for GitHub downloads
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12 -bor [Net.SecurityProtocolType]::Tls13
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $zipPath -UseBasicParsing

    Write-Info "Extracting archive..."
    $extractDir = Join-Path $tempDir "extracted"
    Expand-Archive -Path $zipPath -DestinationPath $extractDir -Force

    # Locate the executable
    $exeSource = Get-ChildItem -Path $extractDir -Filter "*.exe" -Recurse | Select-Object -First 1
    if (-not $exeSource) {
        Write-Err "Could not find an executable inside the downloaded package."
        exit 1
    }

    # 4. Copy to target directory
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    $targetExe = Join-Path $InstallDir "vtablet.exe"
    Copy-Item -Path $exeSource.FullName -Destination $targetExe -Force
    Write-Info "Installed executable to $targetExe"

    # 5. Add to User PATH if requested
    if (-not $NoPathUpdate) {
        $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
        $pathEntries = if ($userPath) { $userPath -split ';' } else { @() }

        if ($pathEntries -notcontains $InstallDir) {
            Write-Info "Adding $InstallDir to user PATH environment variable..."
            $newPath = ($pathEntries + $InstallDir) -join ';'
            [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
            $env:Path = "$env:Path;$InstallDir"
        } else {
            Write-Info "$InstallDir is already in the user PATH."
        }
    }

    # 6. Create Start Menu shortcut
    if (-not $NoShortcut) {
        try {
            $programsDir = [Environment]::GetFolderPath("Programs")
            $shortcutFolder = Join-Path $programsDir "vTablet"
            if (-not (Test-Path $shortcutFolder)) {
                New-Item -ItemType Directory -Path $shortcutFolder -Force | Out-Null
            }

            $shortcutPath = Join-Path $shortcutFolder "vTablet.lnk"
            $wscriptShell = New-Object -ComObject WScript.Shell
            $shortcut = $wscriptShell.CreateShortcut($shortcutPath)
            $shortcut.TargetPath = $targetExe
            $shortcut.WorkingDirectory = $InstallDir
            $shortcut.Description = "vTablet Virtual Digitizer Server"
            $shortcut.Save()

            Write-Info "Created Start Menu shortcut at $shortcutPath"
        } catch {
            Write-Warn "Could not create Start Menu shortcut: $_"
        }
    }

    Write-Info "Installation successful!"
    Write-Info "You can now run 'vtablet' from any terminal, or launch it from the Start Menu."
    Write-Info "For headless mode, run: vtablet -headless"
} finally {
    if (Test-Path $tempDir) {
        Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}
