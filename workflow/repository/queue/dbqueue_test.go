package queue

import (
	"fmt"

	"testing"
	"time"

	"github.com/mmtbak/microlibrary/rdb"

	"gorm.io/gorm"

	"github.com/go-playground/assert/v2"
	"github.com/smartystreets/goconvey/convey"

	"github.com/mmtbak/microlibrary/rdb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSyncSchema(t *testing.T) {

	var err error
	dbclient := getTestDBClient()

	convey.Convey("Test_Connect Database", t, func() {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		convey.So(err, convey.ShouldBeNil)
		dbclient = (&rdb.DBClient{}).WithDB(db)
		convey.So(dbclient, convey.ShouldNotBeNil)

	})

	convey.Convey("Test NewDBDelayMessageQueue", t, func() {
		// DBDelayQueue
		dbopt := DefaultDBDelayQueueOption
		dbmq := NewDBMessageQueue(dbclient, dbopt)
		convey.So(err, convey.ShouldBeNil)
		err = dbmq.SyncSchema()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDelayMessageQueue(t *testing.T) {
	var err error
	dbclient := getTestDBClient()

	// DBDelayQueue
	dbmq := NewDBMessageQueue(dbclient, DefaultDBDelayQueueOption)
	err = dbmq.SyncSchema()
	assert.Equal(t, err, nil)
	fmt.Println("dbmq sync schema success")
	var tables = []string{}
	err = dbclient.DB().Raw("SELECT name FROM sqlite_master WHERE type='table'").Find(&tables).Error
	assert.Equal(t, err, nil)
	fmt.Println("db tables: ", tables)

	rcvchan, err := dbmq.ReceiveInnerMessage()
	assert.Equal(t, err, nil)

	var testcases = []struct {
		Message InnerMessageBody
	}{
		{
			Message: InnerMessageBody{
				ExecutionID: 1,
				Class:       "Execution",
				Type:        "ExecutionInit",
				StepID:      1,
				Data:        `{"testkey":"testvalue"}`,
			},
		},
		{
			Message: InnerMessageBody{
				ExecutionID: 1,
				Class:       "Execution",
				Type:        "ExecutionInit",
				StepID:      2,
				Data:        `{"testkey":"testvalue"}`,
			},
		},
		{
			Message: InnerMessageBody{
				ExecutionID: 2,
				Class:       "Execution",
				Type:        "ExecutionInit",
				StepID:      2,
				Data:        `{"testkey":"testvalue"}`,
			},
		},
		{
			Message: InnerMessageBody{
				ExecutionID: 2,
				Class:       "Execution",
				Type:        "ExecutionInit",
				StepID:      3,
				Data:        `{"testkey":"testvalue"}`,
			},
		},
		{
			Message: InnerMessageBody{
				ExecutionID: 2,
				Class:       "Step",
				Type:        "StepInit",
				StepID:      4,
				Data:        `{"testkey":"testvalue"}`,
			},
		},
		{
			Message: InnerMessageBody{
				ExecutionID: 2,
				Class:       "Execution",
				Type:        "ExecutionInit",
				StepID:      5,
				Data:        `{"testkey":"testvalue"}`,
			},
		},
	}

	rcvchan, err = dbmq.ReceiveInnerMessage()
	assert.Equal(t, err, nil)

	var wait = make(chan int)
	// 接收信息
	go func() {
		var rcvmsg InnerMessageBody
		fmt.Print("\n=> Start Receive Inner Message\n")
		for idx, tt := range testcases {
			rcvmsg = <-rcvchan
			fmt.Printf("\n=> Received index -- %d  Inner Message :  %v  \n %v \n", idx, time.Now(), rcvmsg.Body())
			tt.Message.ID = rcvmsg.Body().ID
			if tt.Message != rcvmsg.Body() {
				fmt.Printf("receive message not match \n")
			} else {
				fmt.Printf("receive message  match \n")
			}
			// Ack inner message
			err = rcvmsg.Ack()
			fmt.Println("ack: ", err)
		}
		fmt.Print("\n=>Receive Inner Message  Finish\n")
		wait <- 0
		fmt.Println("receive finish ")

	}()

	time.Sleep(2 * time.Second)
	for idx, tt := range testcases {
		fmt.Println("index --- ", idx)
		sendtime := time.Now()
		err = dbmq.SendInnerMessage(tt.Message, sendtime)
		assert.Equal(t, err, nil)
		fmt.Println("\n=> Send Normal Inner Message Finish , SendTime: ", time.Now())
	}
	<-wait
	fmt.Println("receive stop ")
	err = dbmq.Close()
	assert.Equal(t, err, nil)
	err = masterinnerqueue.Close()
	assert.Equal(t, err, nil)

	fmt.Printf("finish ")

}
