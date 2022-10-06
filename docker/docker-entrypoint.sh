#!/usr/bin/env sh

set -e

# see: https://unix.stackexchange.com/questions/308260/what-does-set-do-in-this-dockerfile-entrypoint

# if first arg is `-config` or `--some-option`
if [ "${1#-}" != "$1" ]; then
	set -- reciever_ms "$@"
fi

exec "$@"
