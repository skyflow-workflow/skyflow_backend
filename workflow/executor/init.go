package executor

import (
	"github.com/go-playground/validator/v10"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/config"
)

// init validator
var myValidate = validator.New()

// default two executor
var StandardExecutor = NewExecutor(&config.StandardExecutorConfig)
var ExpressExecutor = NewExecutor(&config.ExpressExecutorConfig)
