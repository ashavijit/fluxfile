#!/bin/bash

# Script to install git hooks

HOOK_DIR=".git/hooks"
PRE_PUSH_SCRIPT="scripts/pre-push.sh"
TARGET_HOOK="$HOOK_DIR/pre-push"

if [ ! -d "$HOOK_DIR" ]; then
    echo "Error: .git directory not found. Are you in the root of the repo?"
    exit 1
fi

echo "Installing pre-push hook..."
cp "$PRE_PUSH_SCRIPT" "$TARGET_HOOK"
chmod +x "$TARGET_HOOK"

echo "Hook installed successfully in $TARGET_HOOK"
