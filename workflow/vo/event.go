package vo

import "time"

// ExecutionEvent ...
type ExecutionEvent struct {
	ExecutionID   int
	ExecutionUUID string
	ExecutionURI  string
	StepName      string
	StepID        int
	Data          interface{}
	StartTime     time.Time
	FinishTime    time.Time
}
