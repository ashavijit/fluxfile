$ErrorActionPreference = "Stop"

$OS = "windows"
$ARCH = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

$BIN_URL = "https://github.com/ashavijit/fluxfile/releases/latest/download/flux-${OS}-${ARCH}.exe"
$INSTALL_DIR = "$env:LOCALAPPDATA\flux"
$BIN_PATH = "$INSTALL_DIR\flux.exe"

Write-Host "Downloading Flux for $OS/$ARCH ..."

if (!(Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR | Out-Null
}

Invoke-WebRequest -Uri $BIN_URL -OutFile $BIN_PATH

$USER_PATH = [Environment]::GetEnvironmentVariable("Path", "User")
if ($USER_PATH -notlike "*$INSTALL_DIR*") {
    [Environment]::SetEnvironmentVariable("Path", "$USER_PATH;$INSTALL_DIR", "User")
    Write-Host "Added $INSTALL_DIR to PATH"
}

Write-Host "Flux installed successfully!"
& $BIN_PATH -v
