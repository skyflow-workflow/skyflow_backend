package parser

import (
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/decoder"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/quota"
)

// StandardParserConfig standard model workflow
var StandardParserConfig = decoder.ParserConfig{
	AllowActivity: true,
	AllowWait:     true,
	AllowSuspend:  true,
	AllowParallel: true,
	AllowMap:      true,
	AllowChoice:   true,
	AllowFail:     true,
	AllowSucceed:  true,
	AllowPass:     true,
}

// ExpressParserConfig express model workflow
var ExpressParserConfig = decoder.ParserConfig{
	AllowActivity: false,
	AllowWait:     false,
	AllowSuspend:  false,
	AllowParallel: true,
	AllowMap:      true,
	AllowChoice:   true,
	AllowFail:     true,
	AllowSucceed:  true,
	AllowPass:     true,
}

var StandardParser = NewParser(StandardParserConfig, quota.DefaultQuota)
var ExpressParser = NewParser(ExpressParserConfig, quota.DefaultQuota)
