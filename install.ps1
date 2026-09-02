# ==============================================================================
# SubiKit - Instalador Rápido para Windows (PowerShell)
# Uso: irm https://raw.githubusercontent.com/santi-subidia/dev-kit-desarrollo/main/install.ps1 | iex
# ==============================================================================

$ErrorActionPreference = "Stop"

Write-Host "`n⚡ Instalando SubiKit: Dev-Kit para Desarrollo con IA..." -ForegroundColor Cyan

$owner = "santi-subidia"
$repo = "subi-kit-desarrollo"
$installDir = Join-Path $HOME ".subikit\bin"
$exePath = Join-Path $installDir "subikit.exe"

# 1. Crear directorio de instalación
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

# 2. Obtener la última versión disponible
$tag = "v0.4.0"
try {
    $apiUrl = "https://api.github.com/repos/$owner/$repo/releases/latest"
    $release = Invoke-RestMethod -Uri $apiUrl -Headers @{ "User-Agent" = "SubiKit-Installer" } -ErrorAction SilentlyContinue
    if ($release -and $release.tag_name) {
        $tag = $release.tag_name
    }
} catch {
    # Fallback a release conocido
    $tag = "v0.4.0"
}

# 3. Determinar arquitectura
$arch = "amd64"
if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture -eq [System.Runtime.InteropServices.Architecture]::Arm64) {
    $arch = "arm64"
}

$assetName = "subikit-windows-$arch.zip"
$downloadUrl = "https://github.com/$owner/$repo/releases/download/$tag/$assetName"
$zipPath = Join-Path $env:TEMP "subikit-install-$arch.zip"

Write-Host "-> Descargando SubiKit $tag ($arch)..." -ForegroundColor Gray
try {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing
    
    # 4. Descomprimir binario
    Expand-Archive -Path $zipPath -DestinationPath $installDir -Force
    Remove-Item $zipPath -Force -ErrorAction SilentlyContinue
} catch {
    Write-Host "-> No se pudo descargar el release pre-compilado de GitHub ($downloadUrl)." -ForegroundColor Yellow
    Write-Host "-> Si tienes Go instalado, intentando 'go install' como fallback..." -ForegroundColor Gray
    if (Get-Command go -ErrorAction SilentlyContinue) {
        go install "github.com/$owner/$repo/cmd/subikit@latest"
        Write-Host "✓ SubiKit instalado mediante Go en `$HOME/go/bin" -ForegroundColor Green
        exit 0
    } else {
        Write-Host "[ERROR] No se pudo completar la instalación. Visita https://github.com/$owner/$repo" -ForegroundColor Red
        exit 1
    }
}

# 5. Configurar PATH de usuario si no está presente
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    Write-Host "-> Agregando $installDir al PATH de usuario..." -ForegroundColor Gray
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    $env:Path = "$env:Path;$installDir"
}

Write-Host "`n✓ ¡SubiKit $tag instalado con éxito en $installDir!" -ForegroundColor Green
Write-Host "-> Ejecuta 'subikit tui' para abrir la interfaz interactiva." -ForegroundColor Cyan
Write-Host "-> Ejecuta 'subikit doctor' para verificar tu entorno.`n" -ForegroundColor Cyan
