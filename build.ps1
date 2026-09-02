Write-Host "Compilando Dev-Kit IA (Windows & Linux)..." -ForegroundColor Cyan

New-Item -ItemType Directory -Force -Path "bin" | Out-Null

# Windows
Write-Host "-> Compilando bin/devkit.exe (Windows amd64)..."
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags "-s -w" -o bin/devkit.exe ./cmd/devkit

# Linux
Write-Host "-> Compilando bin/devkit (Linux amd64)..."
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags "-s -w" -o bin/devkit ./cmd/devkit

Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue

Write-Host "`nCompilación exitosa en /bin!" -ForegroundColor Green
Get-ChildItem bin/
