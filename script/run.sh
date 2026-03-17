#!/bin/bash

set -e

echo "╔════════════════════════════════════════════╗"
echo "║      Running CLI Application               ║"
echo "╚════════════════════════════════════════════╝"
echo

cd "$(dirname "$0")/.."

# Build if not already built
if [ ! -f robot-order-system ]; then
    echo "Building application..."
    bash script/build.sh
fi

# Run the application and capture output
echo "Starting robot order system..."
echo

# Run the CLI and redirect output to result.txt
./robot-order-system > result.txt 2>&1

# Also display the output
cat result.txt

if [ $? -eq 0 ]; then
    echo
    echo "✅ CLI execution completed"
    echo "📄 Results saved to: result.txt"
    exit 0
else
    echo
    echo "❌ CLI execution failed"
    exit 1
fi
