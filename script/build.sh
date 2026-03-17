#!/bin/bash

set -e

echo "╔════════════════════════════════════════════╗"
echo "║      Building CLI Application              ║"
echo "╚════════════════════════════════════════════╝"
echo

cd "$(dirname "$0")/.."

# Build the Go application
echo "Building robot-order-system CLI..."
go build -o robot-order-system -v

if [ -f robot-order-system ]; then
    echo
    echo "✅ Build successful"
    echo "📦 Executable: robot-order-system"
    ls -lh robot-order-system
    exit 0
else
    echo
    echo "❌ Build failed"
    exit 1
fi
