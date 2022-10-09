.PHONY: build
build:
	go build -v ./cmd/main/user_server.go

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: run
run: build
	./user_server

.DEFAULT_GOAL := build