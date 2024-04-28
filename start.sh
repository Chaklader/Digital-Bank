#!/bin/sh

set -e

. /app/app.env

#/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
