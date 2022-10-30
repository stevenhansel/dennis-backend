.PHONY: build
build:
	go build -o bin/denji main.go

.PHONY: migration/create
migration/create:
	migrate create -seq -ext=.sql -dir=./database/migrations $(name)
