#!/usr/bin/env bash
set -e

echo "Compilando SubiKit IA (Linux & Windows)..."
mkdir -p bin

# Linux
echo "-> Compilando bin/subikit (Linux amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/subikit ./cmd/subikit

# Windows
echo "-> Compilando bin/subikit.exe (Windows amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/subikit.exe ./cmd/subikit

echo -e "\nCompilación exitosa en /bin!"
ls -la bin/
