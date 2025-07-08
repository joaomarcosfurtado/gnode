#!/bin/bash

set -e

VERSION="v0.1.0"
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH="amd64"

echo "📦 Instalando gnode $VERSION para $OS-$ARCH"

BIN_URL="https://github.com/joaomarcosfurtado/gnode/releases/download/${VERSION}/gnode-${OS}-${ARCH}"

curl -L "$BIN_URL" -o gnode
chmod +x gnode
sudo mv gnode /usr/local/bin/gnode

echo "✅ gnode instalado com sucesso!"
echo "🔁 Você pode rodar: gnode init"
