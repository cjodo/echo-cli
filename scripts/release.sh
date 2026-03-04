#!/bin/bash
set -e

echo "Releasing echo-cli (version from git tag)..."
goreleaser release --clean
