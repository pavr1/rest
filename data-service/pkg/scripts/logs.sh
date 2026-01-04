#!/bin/bash

# View logs from data service containers

cd "$(dirname "$0")/../docker"

if [ -z "$1" ]; then
    docker-compose logs -f
else
    docker-compose logs -f "$1"
fi
