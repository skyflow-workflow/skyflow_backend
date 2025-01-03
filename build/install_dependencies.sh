#!/bin/bash
echo "installing golang dependence tools"
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.60.2
go install github.com/bufbuild/buf/cmd/buf@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/yoheimuta/protolint/cmd/protolint@latest
go install trpc.group/trpc-go/trpc-cmdline/trpc@latest
go install github.com/kyoh86/richgo@latest
trpc setup
