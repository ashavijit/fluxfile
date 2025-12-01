$ErrorActionPreference = "Stop"

function Show-Spinner($Message) {
    $spinner = @('|', '/', '-', '\')
    $i = 0
    while ($global:spin -eq $true) {
        $char = $spinner[$i % $spinner.Length]
        Write-Host -NoNewline "`r[$char] $Message"
        Start-Sleep -Milliseconds 120
        $i++
    }
    Write-Host "`r[OK] $Message"
}

Write-Host ""
Write-Host "== Installing Flux =="
Write-Host ""

$OS = "windows"
$ARCH = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

$BIN_URL = "https://github.com/ashavijit/fluxfile/releases/latest/download/flux-${OS}-${ARCH}.exe"
$INSTALL_DIR = "$env:LOCALAPPDATA\flux"
$BIN_PATH = "$INSTALL_DIR\flux.exe"

if (!(Test-Path $INSTALL_DIR)) {
    Write-Host "[*] Creating installation directory..."
    New-Item -ItemType Directory -Path $INSTALL_DIR | Out-Null
}

$global:spin = $true
Start-Job -ScriptBlock { param($msg) Show-Spinner $msg } -ArgumentList "Downloading Flux..." | Out-Null

try {
    Invoke-WebRequest -Uri $BIN_URL -OutFile $BIN_PATH -ErrorAction Stop
} finally {
    $global:spin = $false
    Start-Sleep -Milliseconds 200
}

$USER_PATH = [Environment]::GetEnvironmentVariable("Path", "User")
if ($USER_PATH -notlike "*$INSTALL_DIR*") {
    Write-Host "[*] Adding flux to PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$USER_PATH;$INSTALL_DIR", "User")
    Write-Host "[OK] PATH updated."
}

Write-Host ""
Write-Host "[OK] Flux installed successfully."
Write-Host "[*] Running flux --version"
Write-Host ""

& $BIN_PATH -v
