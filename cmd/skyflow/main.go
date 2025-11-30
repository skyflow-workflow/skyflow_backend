package main

import (
	"github.com/skyflow-workflow/skyflow_backend/cmd/skyflow/cmd"
	_ "trpc.group/trpc-go/trpc-gateway/plugin/accesslog"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
)

func main() {
	cmd.Execute()
}
