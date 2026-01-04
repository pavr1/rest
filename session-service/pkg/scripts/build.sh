#!/bin/bash
echo "🔨 Building Session Service..."
cd "$(dirname "$0")/.."
go build -o main .
echo "✅ Build complete!"
