#!/bin/bash

# Provider configuration
PROVIDER_NAME="kaizen"
VERSION="0.0.1"
BINARY_NAME="terraform-provider-${PROVIDER_NAME}_v${VERSION}"

# Determine binary extension for Windows
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    BINARY_EXT=".exe"
else
    BINARY_EXT=""
fi

# Build the provider
echo "Building ${BINARY_NAME}${BINARY_EXT}..."
go build -o "${GOPATH}/bin/${BINARY_NAME}${BINARY_EXT}" .

if [ $? -eq 0 ]; then
    echo "✓ Successfully built ${BINARY_NAME}${BINARY_EXT}"
    echo "  Location: ${GOPATH}/bin/${BINARY_NAME}${BINARY_EXT}"
else
    echo "✗ Build failed"
    exit 1
fi
