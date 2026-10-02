#!/bin/sh
set -e

for f in /migrations/*.sql; do
  if [ ! -e "$f" ]; then
    echo "No migrations found"
    exit 0
  fi

  echo "Applying migration: $f"
  psql -v ON_ERROR_STOP=1 \
    --username "$POSTGRES_USER" \
    --dbname "$POSTGRES_DB" \
    -f "$f"
done
