#!/bin/sh
set -eu

export PGPASSWORD="$POSTGRES_PASSWORD"

python3 /seed/clean_data.py /seed/games.csv /tmp/games.clean.csv

psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1 -f /seed/schema.sql
psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1 -c "TRUNCATE TABLE games RESTART IDENTITY;"
psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1 -c "\copy games(original_id, name) FROM '/tmp/games.clean.csv' WITH (FORMAT csv, HEADER true)"
