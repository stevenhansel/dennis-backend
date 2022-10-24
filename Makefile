.PHONY: migration/up
migration/up:
	go run cmd/database/*.go -env $(env) -command up

.PHONY: migration/down
migration/down:
	go run cmd/database/*.go -env $(env) -command down

.PHONY: migration/create
migration/create:
	migrate create -seq -ext=.sql -dir=./database/migrations $(name)
