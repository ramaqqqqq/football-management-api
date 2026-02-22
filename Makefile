run:
	go run cmd/main.go

deps:
	go mod download
	go mod tidy

migrate.up:
	go run migration/main/main.go up

migrate.rollback:
	go run migration/main/main.go rollback

docs-update:
	rm -rf swagger/v1
	$(shell go env GOPATH)/bin/swag init -g cmd/main.go -o swagger/v1 --ot go,json,yaml --pd true

.PHONY: run deps migrate.up migrate.rollback docs-update
