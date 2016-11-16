#!/bin/sh

# /!\ I know it's not a container launch, but no choice with sqlite
rm -f ./test.db;
if ! command -v sqlite3 >/dev/null; then
    echo "Could not sqlite3. You need to install it"
fi
