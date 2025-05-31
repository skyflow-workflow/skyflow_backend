package exporter

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/mmtbak/microlibrary/paging"
	"github.com/skyflow-workflow/skyflow_backbend/mock"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

func TestSendEvents(t *testing.T) {
	var err error
	var myExporter *exporterService
	client := mock.GetMockDBClient()
	dbLis := NewDBListener(client)

	myExporter, err = NewExporterService(dbLis)
	assert.Equal(t, err, nil)
	err = myExporter.SyncSchema()
	assert.Equal(t, err, nil)

	var testEvents = []vo.ExecutionEvent{
		{
			ExecutionID:   65,
			ExecutionUUID: "xxx",
			StateName:     "S1",
			StateID:       0,
			Data: map[string]any{
				"Input":    "{}",
				"Resource": "activity:unitest/add",
			},
			NanoSeconds: time.Now().UnixNano(),
			StartTime:   time.Now().Add(-10 * time.Second),
			FinishTime:  time.Now(),
		},
	}
	myExporter.SendExecutionEvents(testEvents)
}

func TestListExecutionEvents(t *testing.T) {

	var err error
	var myExporter *exporterService
	client := mock.GetMockDBClient()
	dbLis := NewDBListener(client)

	myExporter, err = NewExporterService(dbLis)
	assert.Equal(t, err, nil)
	err = myExporter.SyncSchema()
	assert.Equal(t, err, nil)

	req := vo.ListExecutionEventsRequest{
		ExecutionUUID: "xxx",
		ExecutionID:   19,
		PageRequest:   paging.DefaultPageRequest,
	}

	_, err = myExporter.ListExecutionEvents(req)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, err, nil)

}

func TestListStepEvents(t *testing.T) {
	var err error
	var myExporter *exporterService
	client := mock.GetMockDBClient()
	dbLis := NewDBListener(client)

	myExporter, err = NewExporterService(dbLis)
	assert.Equal(t, err, nil)
	err = myExporter.SyncSchema()
	assert.Equal(t, err, nil)

	req := vo.ListStepEventsRequest{
		StepID:      426,
		PageRequest: paging.DefaultPageRequest,
	}

	_, err = myExporter.ListStepEvents(req)
	if err != nil {
		t.Error(err.Error())
		return
	}
	assert.Equal(t, err, nil)

}
