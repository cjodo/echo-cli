#!/bin/bash
set -e

echo "Building echo-cli..."
go build -buildvcs=true -o bin/echo-cli .
echo "Built: bin/echo-cli"
