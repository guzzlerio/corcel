#!/usr/bin/env bash

# TODO:  This is by no means the preferred way of doing things but it seems one of only a few ways
# which work when trying to get coverage over multiple packages using go test.  This is one
# benefit of tools like ginkgo which has this handling built in and working.

mode="atomic"
set -e
printf "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=${mode} $d
    if [ -f profile.out ]; then
        cat profile.out | tail -n+2 >> coverage.txt
        rm profile.out
    fi
done

sed -i "1 i\mode: ${mode}" coverage.txt

