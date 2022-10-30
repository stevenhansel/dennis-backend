#!/bin/sh
docker run \
        --name certbot \
        --network apinetwork \
        -v /home/ubuntu/nginx/certbot/www:/var/www/certbot:rw \
        -v /home/ubuntu/nginx/certbot/conf:/etc/letsencrypt:rw \
        --rm \
        certbot/certbot:latest \
        certonly \
        --email example@example.com \
        --agree-tos \
        --webroot \
        --webroot-path /var/www/certbot/ \
        -d api.dennis.dog \
