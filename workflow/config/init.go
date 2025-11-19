package config

var StandardExecutorOption = Option{
	AbortOnFail:        false, // standard mode will not fail fast, so user can retry the failed steps.
	AllowActivity:      true,
	AllowWait:          true,
	AllowSuspend:       true,
	AllowParallel:      true,
	AllowMap:           true,
	AllowChoice:        true,
	AllowFail:          true,
	AllowSucceed:       true,
	AllowPass:          true,
	EnableExecuteIndex: true,
	PersistenceStep:    true,
}
var ExpressExecutorOption = Option{
	AllowActivity:      false,
	AllowWait:          false,
	AllowSuspend:       false,
	AllowParallel:      true,
	AllowMap:           true,
	AllowChoice:        true,
	AllowFail:          true,
	AllowSucceed:       true,
	AllowPass:          true,
	EnableExecuteIndex: false,
	PersistenceStep:    false,
	AbortOnFail:        true, // express mode will always fail fast, so no need to set this option.

}

// StandardParserConfig standard model workflow config
var StandardExecutorConfig = Config{
	Option: &StandardExecutorOption,
	Quota:  &DefaultQuota,
}

// ExpressExecutorConfig express model workflow config
var ExpressExecutorConfig = Config{
	Option: &ExpressExecutorOption,
	Quota:  &DefaultQuota,
}
