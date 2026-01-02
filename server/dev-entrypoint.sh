#!/bin/bash
set -e # exit immediately if any command returns an error (non-zero)

# Optional: Only seed if AUTO_SEED environment variable is set (defaults to true)
AUTO_SEED=${AUTO_SEED:-true}

if [ "$AUTO_SEED" = "true" ]; then
    echo "Running database seeder..."
    go run app.go --seed
else
    echo "Skipping database seeding (AUTO_SEED=$AUTO_SEED)"
fi

echo "Starting Air for live reloading..."
exec /go/bin/air -c .air.linux.conf
