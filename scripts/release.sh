#!/bin/bash
set -e

echo "Releasing echo-cli (version from git tag)...$(git tag -l --sort=-creatordate | head -n 1)"
goreleaser release --clean
