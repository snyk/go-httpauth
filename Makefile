SHELL := /bin/bash

GOCMD=go
GOMOD=$(GOCMD) mod
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

all:
	$(info  "completed running make file for golang project")
fmt:
	@go fmt ./...
lint:
	./script/lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	$(GOTEST) ./... -coverprofile cp.out
build:
	$(GOBUILD) -v

build_docker:
	docker build -t gcr.io/snyk-main/<service name>:${CIRCLE_SHA1} .
	docker push gcr.io/snyk-main/<service name>:${CIRCLE_SHA1}

.PHONY: install-req fmt test lint build tidy imports
