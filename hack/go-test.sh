#!/bin/bash
cd $(dirname $0)/../

GO=$(which go)

sudo $GO test ./... --failfast $1 \
| sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
