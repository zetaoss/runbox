#!/bin/bash
if [ $(id -u) -ne 0 ]; then exec sudo bash "$0" "$@"; fi

cd $(dirname $0)/../

GO=$(which go)
set -x
$GO test ./... -v --failfast $1 \
| sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
