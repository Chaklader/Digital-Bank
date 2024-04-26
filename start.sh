#!/bin/sh

set -e

. /app/app.env

echo "start the app"
exec "$@"
