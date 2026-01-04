#!/bin/bash
echo "🔐 Starting Session Service..."
cd "$(dirname "$0")/.."
go run .
