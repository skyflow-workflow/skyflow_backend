SHELL := /bin/bash
GO := go
GOOS ?= $(shell uname -s | tr [:upper:] [:lower:])
GOARCH ?= $(shell go env GOARCH)
GOPATH ?= $(shell go env GOPATH)
PROTOFILE := proto/skyflow.proto
FLAGS := -ldflags="-s -w"

.PHONY: help

target=build

.PHONY: help
help:
	@echo "help                             show usage"
	@echo ""
	@echo "init"
	@echo "install                          install dependency"
	@echo "clean                            clean *.pyc, *.pyo ..."
	@echo "shell                            run neo shell"
	@echo "lint                             lint code"
	@echo ""
	@echo "pip_compile                      lock package versions"
	@echo "test                             run unit test"

.PHONY: clean
clean:
	@echo "cleaning........"
	@rm -rf dist
	@rm -rf build
	@rm -rf .pytest_cache
	@rm -rf .mypy_cache
	@rm -rf .coverage
	@rm -rf .tox
	@rm -rf .eggs
	@rm -rf *.egg-info

.PHONY: shell
shell:
	@echo "running neo shell........"
	@neo

.PHONY: init
init:
	@echo "initializing........"
	@if [ -f go.mod ]; \
	then \
		echo "go.mod already exists"; \
	else \
		${GO} mod init github.com/skyflow-workflow/skyflow_backbend; \
	fi
	${GO} mod tidy

.PHONY: install
install:
	@echo "installing........"
	${SHELL} +x ./install_dependencies.sh

.PHONY: lint
lint:
	@if golangci-lint config path > /dev/null 2>&1; \
	then \
		golangci-lint run --allow-parallel-runners; \
	else \
		echo ".golangci.yml file not found, use golangci-lint with selected linters. Please add .golangci.yml if needed"; \
		golangci-lint run --allow-parallel-runners --disable-all --enable=govet,ineffassign,staticcheck,typecheck; \
	fi

.PHONY: lint_proto
lint_proto:
	@echo "linting proto files........"
	@protolint lint proto/skyflow.proto

.PHONY: pb
pb:
	@echo "generating pb files........"
	@mkdir -p gen/pb gen/apidoc
	@trpc create -p $(PROTOFILE) -o gen/pb --validate=true --lang=go --rpconly --mock=true --nogomod=false
	@trpc apidocs -p $(PROTOFILE) --swagger --swagger-out=gen/apidoc/skyflow.swagger.json
	@trpc apidocs -p $(PROTOFILE) --openapi --openapi-out=gen/apidoc/skyflow.openapi.json

.PHONY: test
test:
	@echo "testing........"
	@go test -v ./...

.PHONY: build
build:
	@echo "building........"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) version && go build