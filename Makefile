.PHONY: build
build:
	go build -v ./cmd/main/main.go

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: run
run: build
	./main

.DEFAULT_GOAL := build