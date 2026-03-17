#!/bin/bash

set -e

echo "╔════════════════════════════════════════════╗"
echo "║        Running Unit Tests                  ║"
echo "╚════════════════════════════════════════════╝"
echo

cd "$(dirname "$0")/.."

# Run go tests with verbose output
go test -v -cover

test_result=$?

if [ $test_result -eq 0 ]; then
    echo
    echo "✅ All tests passed successfully"
    exit 0
else
    echo
    echo "❌ Tests failed with exit code $test_result"
    exit 1
fi
