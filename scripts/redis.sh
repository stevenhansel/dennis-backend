#!/bin/sh

docker run \
  --name redis \
  -p 6379:6379 \
  --volume dennis:/data \
  --rm \
  --detach \
  redis
