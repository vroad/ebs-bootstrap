#!/bin/bash -eu
SOURCE=/go/src/github.com/vroad/ebs-bootstrap
docker run --rm -v "$PWD":$SOURCE -w $SOURCE golang:1.15 \
    bash -c 'make'
