package parser

import (
	"github.com/skyflow-workflow/skyflow_backbend/workflow/config"
)

// StandardParser
var StandardParser *Parser

// ExpressParser ...
var ExpressParser *Parser

func init() {

	StandardParser = NewParser(&config.StandardExecutorConfig)
	ExpressParser = NewParser(&config.ExpressExecutorConfig)
}
