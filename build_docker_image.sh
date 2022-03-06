#!/bin/sh
VERSION=$(git rev-list --abbrev-commit -1 HEAD)
docker build --build-arg VERSION=${VERSION} -f multi-stage.Dockerfile --no-cache --rm -t ripardo/lrw:${VERSION} .
docker image prune -f --filter label=stage=intermediate
docker images -f "dangling=true" -q
