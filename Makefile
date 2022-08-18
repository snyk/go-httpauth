SHELL := /bin/bash

GOCMD=go
GOMOD=$(GOCMD) mod
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

example1: ./cmd/example1/main.go 
	@$(GOBUILD) -o example1 ./cmd/example1/main.go 

build-examples: example1

all:
	$(info  "completed running make file for golang project")
fmt:
	@$(GOCMD) fmt ./...
lint:
	./script/lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	@$(GOTEST) ./... -coverprofile cp.out
build: build-examples
	@$(GOBUILD) -v ./...
clean:
	@$(GOCMD) clean
	@$(GOCMD) clean -testcache
	@rm -f example1

.PHONY: fmt test lint build tidy build-examples clean
