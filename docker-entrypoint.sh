#!/bin/sh
set -eu

# Refresh the cache after host fonts are mounted into the container.
fc-cache -f >/dev/null 2>&1 || true

if ! fc-match -f '%{family}\n' Roboto 2>/dev/null | grep -qi 'Roboto'; then
	echo 'Roboto font not found. Mount a local font directory at /usr/local/share/fonts/custom.' >&2
	exit 1
fi

exec "$@"
