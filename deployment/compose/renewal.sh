#!/bin/sh
docker run \
        --name certbot \
        --network apinetwork \
        -v /home/ubuntu/nginx/certbot/www:/var/www/certbot:rw \
        -v /home/ubuntu/nginx/certbot/conf:/etc/letsencrypt:rw \
        --rm \
        certbot/certbot:latest \
        renew
