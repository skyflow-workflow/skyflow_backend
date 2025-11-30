package queue

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/skyflow-workflow/skyflow_backend/mock"
	"github.com/smartystreets/goconvey/convey"
)

func (dbMQ *DBMessageQueue) CleanUnittestData() error {
	tx := dbMQ.dbClient.DB().Begin()
	tx = tx.Exec("TRUNCATE TABLE message_queues")
	return tx.Error
}

func TestDBMessageQueue_SyncSchema(t *testing.T) {

	var err error
	dbClient := mock.GetMockDBClient()

	convey.Convey("Test Connect to MySQL DB", t, func() {
		// DBDelayQueue
		dbOpt := DefaultDBDelayQueueOption
		dbMQ := NewDBMessageQueue(dbClient, dbOpt)
		convey.So(err, convey.ShouldBeNil)
		convey.Convey("Test SyncSchema", func() {
			err = dbMQ.SyncSchema()
			convey.So(err, convey.ShouldBeNil)

		})
	})
}

func TestDBMessageQueue_WithForwardQueue(t *testing.T) {
	var err error
	dbClient := mock.GetMockDBClient()
	var dbMQ *DBMessageQueue
	// Connect to MySQL DB
	slog.Info("Test Connect to MySQL DB")
	dbOpt := DefaultDBDelayQueueOption
	dbMQ = NewDBMessageQueue(dbClient, dbOpt)
	slog.Info("Test SyncSchema MySQL DB")
	err = dbMQ.SyncSchema()
	assert.Equal(t, err, nil)
	slog.Info("Test Clean Unittest Data")
	err = dbMQ.CleanUnittestData()
	assert.Equal(t, err, nil)

	var messages = []InnerMessageBody{
		{1, 1, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{1, 2, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{2, 2, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{2, 3, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{2, 4, "Step", "StepInit", `{"testkey":"testvalue"}`},
		{2, 5, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
	}

	forwardQueue := NewSimpleInnerQueue()
	dbMQ.SetForwardQueue(forwardQueue)
	dbMQ.StartPolling()

	wg := sync.WaitGroup{}
	wg.Add(1)
	// receive data from forward queue
	go func() {
		slog.Info("=> Receive Inner Message Start")
		rcvChan, err := forwardQueue.ReceiveInnerMessage()
		assert.Equal(t, err, nil)
		msgCount := 0
		for {
			rcvMsg, ok := <-rcvChan
			if !ok {
				break
			}
			msgCount++
			slog.Info("receive message", "msg", rcvMsg)
		}
		slog.Info(fmt.Sprintf("receive message count: %d", msgCount))
		slog.Info("=> Receive Inner Message Finish")
		wg.Done()
		assert.Equal(t, msgCount, len(messages))
	}()

	for _, msg := range messages {
		err = dbMQ.SendInnerMessage(msg, nil)
		assert.Equal(t, err, nil)
		slog.Info("send message", "msg", msg)
	}
	slog.Info(fmt.Sprintf(
		"all messages sent, send message count: %d", len(messages)),
	)
	time.Sleep(2 * time.Second) // wait for messages to be processed
	err = dbMQ.Close()
	assert.Equal(t, err, nil)
	err = forwardQueue.Close()
	assert.Equal(t, err, nil)
	wg.Wait()
}

func TestDBMQSendReceiveMessage(t *testing.T) {
	var err error
	dbClient := mock.GetMockDBClient()
	var dbMQ *DBMessageQueue
	// Connect to MySQL DB
	slog.Info("Test Connect to MySQL DB")
	dbOpt := DefaultDBDelayQueueOption
	dbMQ = NewDBMessageQueue(dbClient, dbOpt)
	slog.Info("Test SyncSchema MySQL DB")
	err = dbMQ.SyncSchema()
	assert.Equal(t, err, nil)
	slog.Info("Test Clean Unittest Data")
	err = dbMQ.CleanUnittestData()
	assert.Equal(t, err, nil)

	var messages = []InnerMessageBody{
		{1, 1, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{1, 2, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{2, 2, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{2, 3, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{2, 4, "Step", "StepInit", `{"testkey":"testvalue"}`},
		{2, 5, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
	}

	dbMQ.StartPolling()

	wg := sync.WaitGroup{}
	wg.Add(1)
	// receive data from forward queue
	go func() {
		slog.Info("=> Start Receive Inner Message from DB Queue")
		rcvChan, err := dbMQ.ReceiveInnerMessage()
		assert.Equal(t, err, nil)
		msgCount := 0
		for {
			rcvMsg, ok := <-rcvChan
			if !ok {
				break
			}
			msgCount++
			slog.Info("receive message", "msg", rcvMsg)
		}
		slog.Info(fmt.Sprintf("receive message count: %d", msgCount))
		slog.Info("=> Receive Inner Message Finish")
		wg.Done()
		assert.Equal(t, msgCount, len(messages))
	}()

	for _, msg := range messages {
		err = dbMQ.SendInnerMessage(msg, nil)
		assert.Equal(t, err, nil)
		slog.Info("send message", "msg", msg)
	}
	slog.Info(fmt.Sprintf(
		"all messages sent, send message count: %d", len(messages)),
	)
	time.Sleep(2 * time.Second) // wait for messages to be processed
	err = dbMQ.Close()
	assert.Equal(t, err, nil)
	wg.Wait()
}
