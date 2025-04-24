package main

import "github.com/skyflow-workflow/skyflow_backbend/server/apiserver"

func main() {

	apiserver := apiserver.NewApiServer()
	apiserver.Start()
}
