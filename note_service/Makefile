.PHONY: build
build:
	go build -v ./cmd/main/note_server.go

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: run
run: build
	./note_server

.DEFAULT_GOAL := build