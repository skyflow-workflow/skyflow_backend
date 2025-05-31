package vo

import "time"

// ExecutionEvent execution event
type ExecutionEvent struct {
	ID            int
	ExecutionID   int
	ExecutionUUID string
	ExecutionURI  string
	StateName     string
	StateID       int
	EventType     string
	NanoSeconds   int64
	Data          any
	StartTime     time.Time
	FinishTime    time.Time
}
