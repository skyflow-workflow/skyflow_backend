package queue

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/smartystreets/goconvey/convey"
)

func (dbmq *DBMessageQueue) CleanUnittestData() error {
	tx := dbmq.dbClient.DB().Begin()
	tx = tx.Exec("TRUNCATE TABLE message_queues")
	return tx.Error
}

func TestDBMQSyncSchema(t *testing.T) {

	var err error
	dbclient := getTestDBClient()

	convey.Convey("Test Connect to MySQL DB", t, func() {
		// DBDelayQueue
		dbopt := DefaultDBDelayQueueOption
		dbmq := NewDBMessageQueue(dbclient, dbopt)
		convey.So(err, convey.ShouldBeNil)
		convey.Convey("Test SyncSchema", func() {
			err = dbmq.SyncSchema()
			convey.So(err, convey.ShouldBeNil)

		})
	})
}

func TestDBMQWithForwardQueue(t *testing.T) {
	var err error
	dbclient := getTestDBClient()
	var dbmq *DBMessageQueue
	// Connect to MySQL DB
	slog.Info("Test Connect to MySQL DB")
	dbopt := DefaultDBDelayQueueOption
	dbmq = NewDBMessageQueue(dbclient, dbopt)
	slog.Info("Test SyncSchema MySQL DB")
	err = dbmq.SyncSchema()
	assert.Equal(t, err, nil)
	slog.Info("Test Clean Unittest Data")
	err = dbmq.CleanUnittestData()
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
	dbmq.SetForwardQueue(forwardQueue)
	dbmq.StartPolling()

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
		err = dbmq.SendInnerMessage(msg, nil)
		assert.Equal(t, err, nil)
		slog.Info("send message", "msg", msg)
	}
	slog.Info(fmt.Sprintf(
		"all messages sent, send message count: %d", len(messages)),
	)
	time.Sleep(2 * time.Second) // wait for messages to be processed
	err = dbmq.Close()
	assert.Equal(t, err, nil)
	err = forwardQueue.Close()
	assert.Equal(t, err, nil)
	wg.Wait()
}

func TestDBMQSendReceiveMessage(t *testing.T) {
	var err error
	dbclient := getTestDBClient()
	var dbmq *DBMessageQueue
	// Connect to MySQL DB
	slog.Info("Test Connect to MySQL DB")
	dbopt := DefaultDBDelayQueueOption
	dbmq = NewDBMessageQueue(dbclient, dbopt)
	slog.Info("Test SyncSchema MySQL DB")
	err = dbmq.SyncSchema()
	assert.Equal(t, err, nil)
	slog.Info("Test Clean Unittest Data")
	err = dbmq.CleanUnittestData()
	assert.Equal(t, err, nil)

	var messages = []InnerMessageBody{
		{1, 1, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{1, 2, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{2, 2, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{2, 3, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
		{2, 4, "Step", "StepInit", `{"testkey":"testvalue"}`},
		{2, 5, "Execution", "ExecutionInit", `{"testkey":"testvalue"}`},
	}

	dbmq.StartPolling()

	wg := sync.WaitGroup{}
	wg.Add(1)
	// receive data from forward queue
	go func() {
		slog.Info("=> Start Receive Inner Message from DB Queue")
		rcvChan, err := dbmq.ReceiveInnerMessage()
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
		err = dbmq.SendInnerMessage(msg, nil)
		assert.Equal(t, err, nil)
		slog.Info("send message", "msg", msg)
	}
	slog.Info(fmt.Sprintf(
		"all messages sent, send message count: %d", len(messages)),
	)
	time.Sleep(2 * time.Second) // wait for messages to be processed
	err = dbmq.Close()
	assert.Equal(t, err, nil)
	wg.Wait()
}
