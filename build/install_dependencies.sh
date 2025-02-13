#!/bin/bash
echo "installing golang dependence tools"
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.60.2
go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/yoheimuta/protolint/cmd/protolint@latest
go install github.com/kyoh86/richgo@latest
# install trpc tools
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install golang.org/x/tools/cmd/goimports@latest
go install go.uber.org/mock/mockgen@latest
go install github.com/envoyproxy/protoc-gen-validate@latest
go install github.com/envoyproxy/protoc-gen-validate/cmd/protoc-gen-validate-go@latest
go install trpc.group/trpc-go/trpc-cmdline/trpc@latest
trpc setup
