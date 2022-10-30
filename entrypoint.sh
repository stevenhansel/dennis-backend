#!/bin/bash -e

APP_ENV=${APP_ENV:-production}

echo "[`date`] Running entrypoint script in the '${APP_ENV}' environment..."
echo "[`date`] Starting server on http://localhost:${PORT}"

./denji api run --env ${APP_ENV}