.PHONY: up
up:
	docker compose --env-file .env up -d

.PHONY: down
down:
	docker compose down

.PHONY: restart
restart:
	docker compose stop $(container) && docker compose --env-file .env up -d $(container)

.PHONY: rebuild
rebuild:
	docker compose stop $(container) && docker compose --env-file .env up -d --force-recreate --no-deps --build $(container)
