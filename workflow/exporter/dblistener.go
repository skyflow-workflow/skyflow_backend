package exporter

import (
	"log/slog"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// DBListener db as listener
type DBListener struct {
	dbClient *rdb.DBClient
}

func NewDBListener(client *rdb.DBClient) *DBListener {

	return &DBListener{
		dbClient: client,
	}
}

func (l *DBListener) Client() *rdb.DBClient {
	return l.dbClient
}

// SyncSchema implements Listener.
func (l *DBListener) SyncSchema() error {
	err := l.dbClient.SyncTables(po.GetEventTables())
	return err
}

// SendEvents implements Listener Interface.
func (l *DBListener) SendEvents(events []vo.ExecutionEvent) {

	var err error
	var dbEvents []po.ExecutionEvent
	for _, evt := range events {
		dataStr, err := JSONString(evt.Data)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		dbEvent := po.ExecutionEvent{
			ExecutionID: evt.ExecutionID,
			StateID:     evt.StateID,
			StateName:   evt.StateName,
			EventType:   evt.EventType,
			Data:        dataStr,
			NanoSeconds: evt.NanoSeconds,
			StartTime:   evt.StartTime,
			FinishTime:  evt.FinishTime,
		}

		dbEvents = append(dbEvents, dbEvent)
	}
	err = l.Client().DB().Create(dbEvents).Error
	if err != nil {
		slog.Error(err.Error())
	}
}
