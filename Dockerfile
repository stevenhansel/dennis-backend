# Stage 0: building the binary
FROM golang:alpine AS build

RUN apk update && \
    apk add curl \
            git \
            bash \
            make \
            ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.* ./
RUN go mod download
RUN go mod verify

COPY . .
RUN make build

# Stage 1: copying files
FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /app/

COPY --from=build /app/bin/denji .
COPY --from=build /app/database/migrations ./database/migrations/
COPY --from=build /app/entrypoint.sh .

RUN ls -la

ENTRYPOINT ["./entrypoint.sh"]
