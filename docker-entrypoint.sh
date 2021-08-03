#!/bin/sh

# Abort on any error (including if wait-for-it fails).
set -e
set -u

#!/bin/sh
echo "starting container $hostname"
( : ${DATABASE_HOST?"pls pass -e DATABASE_HOST=somehost"} ) || exit 1

HOST=$1
PORT=$2

# Wait for the backend to be up, if we know where it is.
if [ -n "$CUSTOMERS_HOST" ]; then
  /wait-for-it.sh "$CUSTOMERS_HOST:${CUSTOMERS_PORT:-6000}"
fi

# Run the main container command.
exec "$@"

