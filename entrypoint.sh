#!/bin/sh
# Entrypoint script for LeakyRepo Docker Action
# This script handles arguments from GitHub Actions and passes them to leakyrepo

set -e

# Fix Git "dubious ownership" error in GitHub Actions
# The workspace is owned by root but we need to access it
# Handle both /github/workspace (GitHub Actions) and /workspace (custom workflows)
if command -v git >/dev/null 2>&1; then
    if [ -d "/github/workspace" ]; then
    git config --global --add safe.directory /github/workspace 2>/dev/null || true
    fi
    if [ -d "/workspace" ]; then
        git config --global --add safe.directory /workspace 2>/dev/null || true
    fi
    # Also add current working directory if it's a git repo
    if [ -d ".git" ]; then
        git config --global --add safe.directory "$(pwd)" 2>/dev/null || true
    fi
fi

# If no arguments provided, use the default CMD from Dockerfile
if [ $# -eq 0 ]; then
    exec /usr/local/bin/leakyrepo scan --help
    exit 0
fi

# If first argument contains spaces, it's likely a string from GitHub Actions
# Split it into separate arguments
if [ $# -eq 1 ] && echo "$1" | grep -q ' '; then
    # Use eval to properly split the arguments
    eval "set -- $1"
fi

# Execute leakyrepo with all arguments
exec /usr/local/bin/leakyrepo "$@"

