#!/usr/bin/env bash
# Download starlark-unified-schema from GitHub and extract to .cache/
# Usage: ./scripts/fetch-schema.sh [commit-hash-or-branch]

set -euo pipefail

REPO="project-kessel/starlark-unified-schema"
REF="${1:-main}"
CACHE_DIR=".cache/starlark-unified-schema"

# If cache already exists and no specific ref requested, skip download
if [ -d "$CACHE_DIR" ] && [ "$REF" = "main" ]; then
    echo "Schema cache exists at $CACHE_DIR (use 'make refresh-schema' to update)"
    exit 0
fi

echo "Downloading starlark-unified-schema@${REF} from GitHub..."

# Clean and create cache directory
rm -rf "$CACHE_DIR"
mkdir -p "$(dirname "$CACHE_DIR")"

# Download tarball from GitHub and extract (strip top-level dir)
TARBALL_URL="https://github.com/${REPO}/archive/${REF}.tar.gz"
curl -fsSL "$TARBALL_URL" | tar xz -C "$(dirname "$CACHE_DIR")"

# Rename extracted directory to match expected name
EXTRACTED_DIR=$(dirname "$CACHE_DIR")/starlark-unified-schema-${REF}
if [ -d "$EXTRACTED_DIR" ]; then
    mv "$EXTRACTED_DIR" "$CACHE_DIR"
else
    # If ref is a commit hash, GitHub uses a different naming
    POSSIBLE_DIR=$(dirname "$CACHE_DIR")/starlark-unified-schema-*
    if compgen -G "$POSSIBLE_DIR" > /dev/null; then
        mv $POSSIBLE_DIR "$CACHE_DIR"
    else
        echo "Error: Could not find extracted directory" >&2
        exit 1
    fi
fi

# Verify required directories exist
if [ ! -d "$CACHE_DIR/schema" ]; then
    echo "Error: schema/ directory not found in downloaded archive" >&2
    exit 1
fi

if [ ! -d "$CACHE_DIR/interpreter" ]; then
    echo "Error: interpreter/ directory not found in downloaded archive" >&2
    exit 1
fi

echo "Repository cached to $CACHE_DIR"
echo "  - Go module: $CACHE_DIR/interpreter"
echo "  - Schema files: $CACHE_DIR/schema"
