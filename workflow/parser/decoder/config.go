package decoder

// ParserConfig ...
type ParserConfig struct {
	// AllowActivity specifies whether to allow activity Task.
	AllowActivity bool
	// AllowWait specifies whether to allow Wait State.
	AllowWait bool
	// AllowSuspend specifies whether to allow Suspend State.
	AllowSuspend bool
	// AllowParallel specifies whether to allow Parallel State.
	AllowParallel bool
	// AllowMap specifies whether to allow Map State.
	AllowMap bool
	// AllowChoice specifies whether to allow Choice State.
	AllowChoice bool
	// AllowFail specifies whether to allow Fail State.
	AllowFail bool
	// AllowSucceed specifies whether to allow Succeed State.
	AllowSucceed bool
	// AllowPass specifies whether to allow Pass State.
	AllowPass bool
}

// StandardParserConfig standard model workflow
var StandardParserConfig = ParserConfig{
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
