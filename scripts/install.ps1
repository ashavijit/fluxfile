$ErrorActionPreference = "Stop"

# Colors
$Cyan = [char]27 + "[36m"
$Green = [char]27 + "[32m"
$Yellow = [char]27 + "[33m"
$Red = [char]27 + "[31m"
$Reset = [char]27 + "[0m"
$Bold = [char]27 + "[1m"

function Write-Step {
    param($Icon, $Message, $Color = $Reset)
    Write-Host "  $Color$Icon$Reset $Message"
}

function Write-Header {
    Write-Host ""
    Write-Host "  $Cyan${Bold}=======================================$Reset"
    Write-Host "  $Cyan${Bold}          FLUX INSTALLER               $Reset"
    Write-Host "  $Cyan${Bold}=======================================$Reset"
    Write-Host ""
}

Write-Header

$OS = "windows"
$ARCH = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$BIN_URL = "https://github.com/ashavijit/fluxfile/releases/latest/download/flux-${OS}-${ARCH}.exe"
$INSTALL_DIR = "$env:LOCALAPPDATA\flux"
$BIN_PATH = "$INSTALL_DIR\flux.exe"

# Create install directory
if (!(Test-Path $INSTALL_DIR)) {
    Write-Step "[*]" "Creating installation directory..."
    New-Item -ItemType Directory -Path $INSTALL_DIR | Out-Null
}

# Remove old binary if exists
if (Test-Path $BIN_PATH) {
    Write-Step "[!]" "Removing old version..." $Yellow
    try {
        Remove-Item $BIN_PATH -Force -ErrorAction Stop
        Write-Step "[OK]" "Old version removed" $Green
    } catch {
        Write-Step "[!]" "Could not remove old version (may be in use)" $Yellow
    }
}

# Download new binary
Write-Step "[*]" "Downloading Flux ($OS-$ARCH)..."
try {
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest -Uri $BIN_URL -OutFile $BIN_PATH -ErrorAction Stop
    Write-Step "[OK]" "Download complete" $Green
} catch {
    Write-Step "[X]" "Download failed: $_" $Red
    exit 1
}

# Update PATH
$USER_PATH = [Environment]::GetEnvironmentVariable("Path", "User")
if ($USER_PATH -notlike "*$INSTALL_DIR*") {
    Write-Step "[*]" "Adding flux to PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$USER_PATH;$INSTALL_DIR", "User")
    Write-Step "[OK]" "PATH updated" $Green
} else {
    Write-Step "[OK]" "PATH already configured" $Green
}

# Verify installation
Write-Host ""
Write-Host "  $Green${Bold}=======================================$Reset"
Write-Host "  $Green${Bold}       INSTALLATION COMPLETE           $Reset"
Write-Host "  $Green${Bold}=======================================$Reset"
Write-Host ""

Write-Step "[>]" "Installed to: $BIN_PATH"

try {
    $version = & $BIN_PATH -v 2>&1
    Write-Step "[>]" "Version: $version"
} catch {
    Write-Step "[!]" "Could not verify version" $Yellow
}

Write-Host ""
Write-Host "  ${Cyan}Usage:$Reset"
Write-Host "    flux init            Create new FluxFile"
Write-Host "    flux build           Run build task"
Write-Host "    flux -l              List all tasks"
Write-Host "    flux logs            View execution logs"
Write-Host ""
