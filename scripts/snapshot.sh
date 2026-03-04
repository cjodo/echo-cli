#!/bin/bash
set -e

echo "Creating snapshot release..."
goreleaser release --clean --snapshot --skip=publish
