#!/bin/sh

docker run \
  --name postgres \
  --network chainsawman \
  --volume chainsawman:/var/lib/postgresql/data \
  -p 5432:5432 \
  --rm \
  --detach \
  -e POSTGRES_USER=denji \
  -e POSTGRES_PASSWORD=c4c6ce446bd2a17a4764b9aada51e5d243471512643ec84d43d5fccabb5cb2d0 \
  -e POSTGRES_DB=chainsawman \
  postgres
