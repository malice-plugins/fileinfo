#!/bin/bash
set -e

# Add info as command if needed
if [[ "$1" == -* ]]; then
	set -- info "$@"
fi

# Run as user "malice" if the command is "info"
if [ "$1" = 'info' ]; then
	set -- gosu malice tini -- "$@"
fi

# else default to run whatever the user wanted like "bash"
exec "$@"
