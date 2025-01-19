#!/bin/sh

set -e

HOST="\$1"
PORT="\$2"

echo "Trying to connect to Postgres at $HOST:$PORT"

until pg_isready -h "$HOST" -p "$PORT"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec "${@:3}"
