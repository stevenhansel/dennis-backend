#!/bin/sh

docker run \
  --name postgres \
  --network chainsawman \
  --volume chainsawman:/var/lib/postgresql/data \
  -p 5432:5432 \
  --rm \
  --detach \
  -e POSTGRES_USER=denji \
  -e POSTGRES_PASSWORD=woofwoof \
  -e POSTGRES_DB=chainsawman \
  postgres
