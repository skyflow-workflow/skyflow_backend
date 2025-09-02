package executor

import (
	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/exporter"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

type ExecutionService = *executionService

type executionService struct {
	MetaDB     *rdb.DBClient
	InnerQueue queue.InnerMessageQueue
	Exporter   exporter.ExporterService
}

func NewExecutionService(
	MetaDB *rdb.DBClient,
	InnerQueue queue.InnerMessageQueue,
	Exporter exporter.ExporterService,
) ExecutionService {
	return &executionService{
		MetaDB:     MetaDB,
		InnerQueue: InnerQueue,
		Exporter:   Exporter,
	}
}

// SendExecutionEvents 发送event
func (svc *executionService) SendExecutionEvents(events ...vo.ExecutionEvent) {
	if svc.Exporter == nil {
		return
	}
	svc.Exporter.SendExecutionEvents(events)
}
