#!/bin/bash

# Module Rename Script
# Usage: ./module.sh [old_path] [new_path]

set -euo pipefail

# Default values
OLD_PATH="${1:-webservices}"
NEW_PATH="${2:-webservices}"

echo "Renaming Go module from $OLD_PATH to $NEW_PATH"

# 1. Update go.mod
echo "Updating go.mod..."
if [ ! -f go.mod ]; then
    echo "Error: go.mod file not found!"
    exit 1
fi

sed -i "s|$OLD_PATH|$NEW_PATH|g" go.mod

# 2. Update all .go files
echo "Updating imports in Go files..."
find . -name "*.go" -type f | while read -r file; do
    sed -i "s|\"$OLD_PATH|\"$NEW_PATH|g" "$file"
done

# 3. Clean and verify
echo "Running go mod tidy..."
go mod tidy

# echo "Verifying build..."
# go build ./...

echo "Module rename completed successfully!"
echo "Changed from: $OLD_PATH"
echo "Changed to:   $NEW_PATH"