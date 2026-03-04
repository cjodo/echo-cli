#!/bin/bash
set -e

echo "Running echo-cli in dev mode..."
go run -buildvcs=true . "$@"
