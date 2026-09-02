#!/usr/bin/env bash
set -e

echo "Compilando Dev-Kit IA (Linux & Windows)..."
mkdir -p bin

# Linux
echo "-> Compilando bin/devkit (Linux amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/devkit ./cmd/devkit

# Windows
echo "-> Compilando bin/devkit.exe (Windows amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/devkit.exe ./cmd/devkit

echo -e "\nCompilación exitosa en /bin!"
ls -la bin/
