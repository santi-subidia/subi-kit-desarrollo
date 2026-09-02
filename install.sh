#!/usr/bin/env bash
# ==============================================================================
# SubiKit - Instalador Rápido para Linux y macOS
# Uso: curl -fsSL https://raw.githubusercontent.com/santi-subidia/dev-kit-desarrollo/main/install.sh | bash
# ==============================================================================

set -e

OWNER="santi-subidia"
REPO="subi-kit-desarrollo"
INSTALL_DIR="$HOME/.local/bin"

echo -e "\n\033[36m⚡ Instalando SubiKit: Dev-Kit para Desarrollo con IA...\033[0m"

# 1. Detectar Sistema Operativo y Arquitectura
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "\033[31m[ERROR] Arquitectura no soportada: $ARCH\033[0m"
        exit 1
        ;;
esac

if [ "$OS" != "linux" ] && [ "$OS" != "darwin" ]; then
    echo -e "\033[31m[ERROR] Sistema operativo no soportado: $OS\033[0m"
    exit 1
fi

mkdir -p "$INSTALL_DIR"

# 2. Obtener versión más reciente
TAG="v0.4.0"
API_URL="https://api.github.com/repos/$OWNER/$REPO/releases/latest"
LATEST_TAG=$(curl -sSL "$API_URL" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || true)
if [ -n "$LATEST_TAG" ]; then
    TAG="$LATEST_TAG"
fi

ASSET_NAME="subikit-${OS}-${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$OWNER/$REPO/releases/download/$TAG/$ASSET_NAME"

echo -e "-> Descargando SubiKit $TAG ($OS/$ARCH)..."

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

if curl -fL "$DOWNLOAD_URL" -o "$TMP_DIR/$ASSET_NAME" 2>/dev/null; then
    tar -xzf "$TMP_DIR/$ASSET_NAME" -C "$TMP_DIR"
    chmod +x "$TMP_DIR/subikit"
    mv "$TMP_DIR/subikit" "$INSTALL_DIR/subikit"
else
    echo -e "\033[33m-> No se pudo descargar el release pre-compilado de GitHub ($DOWNLOAD_URL).\033[0m"
    if command -v go >/dev/null 2>&1; then
        echo "-> Compilando e instalando con Go..."
        go install "github.com/$OWNER/$REPO/cmd/subikit@latest"
        echo -e "\033[32m✓ SubiKit instalado mediante Go en $(go env GOPATH)/bin\033[0m"
        exit 0
    else
        echo -e "\033[31m[ERROR] No se pudo completar la instalación. Visita https://github.com/$OWNER/$REPO\033[0m"
        exit 1
    fi
fi

# 3. Comprobar PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "\n\033[33mAviso: $INSTALL_DIR no está en tu PATH.\033[0m"
    echo "Agrega la siguiente línea a tu ~/.bashrc o ~/.zshrc:"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

echo -e "\n\033[32m✓ ¡SubiKit $TAG instalado con éxito en $INSTALL_DIR/subikit!\033[0m"
echo -e "\033[36m-> Ejecuta 'subikit tui' para abrir la interfaz interactiva.\033[0m"
echo -e "\033[36m-> Ejecuta 'subikit doctor' para verificar tu entorno.\n\033[0m"
