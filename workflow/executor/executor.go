package executor

import (
	"time"
)

type ExecutorConfig struct {
	Timeout time.Duration
}

// StandardExecutorConfig is the configuration for the standard executor
var StandardExecutorConfig = ExecutorConfig{
	Timeout: 30 * 24 * time.Hour,
}

// ExpressExecutorConfig is the configuration for the express executor
var ExpressExecutorConfig = ExecutorConfig{
	Timeout: 5 * time.Minute,
}

// SyncExecutorConfig is the configuration for the sync executor
// The sync executor is used for the sync workflow to orchestrate microservices
var SyncExecutorConfig = ExecutorConfig{
	Timeout: 1 * time.Minute,
}
