#!/bin/sh

echo "Setting up Git hooks..."

# Ensure .git/hooks exists
mkdir -p .git/hooks

# Link the pre-commit hook
ln -sf ../../.husky/pre-commit .git/hooks/pre-commit


# make the hook executable
chmod +x .husky/pre-commit


echo "Git hooks setup completed."