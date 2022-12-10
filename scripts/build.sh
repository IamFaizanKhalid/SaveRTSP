#!/usr/bin/env bash

source ./scripts/env.sh

if [ ! -d $OUTPUT_DIRECTORY ]; then
    echo "Creating Directory \`$OUTPUT_DIRECTORY\`.."
    mkdir -p $OUTPUT_DIRECTORY
fi

go build -o $BINARY
