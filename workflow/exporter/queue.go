package exporter

import (
	"log/slog"

	"github.com/mmtbak/microlibrary/mq"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

var (
	ContentType_JSON          = "application/json"
	ContentType_CloudEvent2_0 = "cloudevent2.0"
)

type QueueListener struct {
	queue       mq.MessageQueue
	ContentType string
}

func NewQueueListener(queue mq.MessageQueue) *QueueListener {
	q := &QueueListener{
		queue:       queue,
		ContentType: ContentType_JSON,
	}
	return q
}

// SyncSchema implements Listener.
func (q *QueueListener) SyncSchema() error {
	return q.queue.SyncSchema()
}

// SendEvents implements Listener Interface.
func (q *QueueListener) SendEvents(events []vo.ExecutionEvent) {

	for _, evt := range events {
		data, err := JSONBytes(evt.Data)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		q.queue.SendMessage(data)
	}

}
