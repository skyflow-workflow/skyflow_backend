SHELL := /bin/bash
GO := go
GOOS ?= $(shell uname -s | tr [:upper:] [:lower:])
# GOTEST ?= ${GO} test
GOTEST ?= richgo test
GOARCH ?= $(shell go env GOARCH)
GOPATH ?= $(shell go env GOPATH)
GOPROJECT := github.com/skyflow-workflow/skyflow_backbend
PROTOFILE := proto/skyflow.proto
FLAGS := -ldflags="-s -w"

.PHONY: help

all: build

.PHONY: help
help:
	@echo "help                             show usage"
	@echo "init	                            initialize go mod"
	@echo "install                          install dependency"
	@echo "clean                            clean tmp files"
	@echo "lint                             lint code"
	@echo "lint_proto                       lint proto files"
	@echo "test                             run unittest"

.PHONY: clean
clean:
	@echo "cleaning........"
	@rm -rf dist

.PHONY: init
init:
	@echo "initializing........"
	@if [ -f go.mod ]; \
	then \
		echo "go.mod already exists"; \
	else \
		${GO} mod init ${GOPROJECT}; \
	fi
	${GO} mod tidy

.PHONY: install
install:
	@echo "installing........"
	${SHELL} +x ./build/install_dependencies.sh

.PHONY: lint
lint:
	golangci-lint run --allow-parallel-runners

.PHONY: lint_proto
lint_proto:
	@echo "linting proto files........"
	@protolint lint ${PROTOFILE}

.PHONY: pb
pb:
	@echo "generating pb files........"
	@mkdir -p gen/pb gen/apidoc
	@trpc create -p $(PROTOFILE) -o gen/pb --validate=true --protocol=trpc --lang=go --rpconly --mock=true --nogomod=false
	# @trpc create -p $(PROTOFILE) -o gen/pb --validate=true --protocol=trpc --lang=go --rpconly --mock=true --nogomod=true
	@trpc apidocs -p $(PROTOFILE) --swagger --swagger-out=gen/apidoc/skyflow.swagger.json
	@trpc apidocs -p $(PROTOFILE) --openapi --swagger=false --openapi-out=gen/apidoc/skyflow.openapi.json

.PHONY: test
test:
	@echo "start unittest........"
	${GOTEST} -v ./...

.PHONY: build
build:
	@echo "building........"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) version && go build  -o bin/skyflow $(FLAGS) ./cmd/skyflow/*.go

.PHONY: run
run: build
	@echo "running........"
	@./bin/skyflow -conf trpc_go.yaml