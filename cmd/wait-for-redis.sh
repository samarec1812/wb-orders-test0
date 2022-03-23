#!/usr/bin/env bash

set -e

HOST=$1
PORT=$2


echo "Start waiting for Redis fully start. Host '$HOST', '$PORT'..."
echo "Try ping Redis... "
PONG=`redis-cli -h $HOST -p $PORT ping | grep PONG`
while [ -z "$PONG" ]; do
    sleep 1
    >&2 echo "Retry Redis ping... "
    PONG=`redis-cli -h $HOST -p $PORT ping | grep PONG`
done
>&2 echo "Redis at host '$HOST', port '$PORT' fully started."
shift
shift
exec "$@"