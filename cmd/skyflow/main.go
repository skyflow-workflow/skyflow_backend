package main

import (
	"github.com/skyflow-workflow/skyflow_backbend/cmd/skyflow/cmd"
	_ "trpc.group/trpc-go/trpc-gateway/plugin/accesslog"
)

func main() {
	cmd.Execute()
}
