#!/bin/bash
set -e

echo "Running DB migrations for testing..."
migrate -path migrations -database "${DATABASE_URL}" up

echo "Executing command: $@"
exec "$@"
