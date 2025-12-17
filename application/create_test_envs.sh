#!/bin/bash

# Script to create test environment files from .env.test
# Creates .env.test_complete, .env.test_normal, .env.test_small
# and modifies DBNAME parameter to match the suffix

set -e

SOURCE_FILE=".env.test"
TARGETS=("test_complete" "test_normal" "test_small")

# Check if source file exists
if [ ! -f "$SOURCE_FILE" ]; then
    echo "Error: $SOURCE_FILE not found"
    exit 1
fi

echo "Creating test environment files from $SOURCE_FILE..."

for target in "${TARGETS[@]}"; do
    dest_file=".env.$target"

    # Copy source file to destination
    cp "$SOURCE_FILE" "$dest_file"

    # Replace DBNAME value with the target name
    sed -i "s/^DBNAME=.*/DBNAME=$target/" "$dest_file"

    echo "Created $dest_file with DBNAME=$target"
done

echo "Done! Created ${#TARGETS[@]} environment files."
