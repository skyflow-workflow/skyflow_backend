package exporter

import "github.com/skyflow-workflow/skyflow_backbend/workflow/vo"

// Listener 日志监听器
type Listener interface {
	SendEvents([]vo.ExecutionEvent)
	SyncSchema() error
}
