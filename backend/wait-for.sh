#!/bin/sh

HOST=$(echo "$1" | cut -d: -f1)
PORT=$(echo "$1" | cut -d: -f2)
shift

echo "Waiting for $HOST:$PORT to be available..."

while ! nc -z "$HOST" "$PORT"; do
  sleep 1
done

echo "$HOST:$PORT is available. Running command..."
exec "$@"