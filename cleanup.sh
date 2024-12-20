#!/bin/bash

set -e

# Check if README.md exists
if [ ! -f "README.md" ]; then
  error "README.md not found in the current directory."
fi

# Configure Git
git config --global --add safe.directory /github/workspace

# Configure Git user
git config --global user.name "github-actions[bot]"
git config --global user.email "github-actions[bot]@users.noreply.github.com"

# Add README.md to staging
git add README.md

# Check if there are any changes to commit
if git diff --cached --quiet; then
  echo "No changes detected in README.md. Skipping commit."
  exit 0
fi

# Commit the changes with the specified message
git commit -m "chore: update readme stats"
git push

echo "Pushed up changes to README.md"
