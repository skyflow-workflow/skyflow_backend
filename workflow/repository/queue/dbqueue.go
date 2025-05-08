package queue

/*
 * @Author: mumangtao@gmail.com
 * @Date: 2020-07-24 20:25:30
 * @Last Modified by: mumangtao@gmail.com
 * @Last Modified time: 2020-07-24 23:50:24
 */

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/mmtbak/microlibrary/config"
	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
)

// For Simulate MessageQueue By DataBase like MySQL

// DBQueueOption db 队列的默认参数
type DBQueueOption struct {
	// 正常消息的巡检时长
	PollingDuration time.Duration
	// 正常消息每次捞出的消息数量
	PollingLimit int
}

// 默认DB队列参数
var DefaultDBDelayQueueOption = DBQueueOption{
	// 每500ms从DB中查找一次消息
	PollingDuration: 500 * time.Millisecond,
	// 每次最多查找10条未读消息
	PollingLimit: 10,
}

func NewDBQueueOption() DBQueueOption {
	return DefaultDBDelayQueueOption
}

// DBMessageStatus message status in db
var DBMessageStatus = struct {
	Created   string
	Processed string
}{
	Created:   "Created",
	Processed: "Processed",
}

// DBDelayMessageQueue DB delay message queue
type DBMessageQueue struct {
	dbclient *rdb.DBClient
	wg       sync.WaitGroup
	option   DBQueueOption
	logger   *slog.Logger
	ctx      context.Context
	cancel   context.CancelFunc
	msgchan  chan InnerMessage
}

// NewDBMessageQueue Create New Inner DBMessage
func NewDBMessageQueue(client *rdb.DBClient, option DBQueueOption) *DBMessageQueue {

	ctx, cancel := context.WithCancel(context.Background())

	dbmq := &DBMessageQueue{
		dbclient: client,
		wg:       sync.WaitGroup{},
		ctx:      ctx,
		cancel:   cancel,
		option:   option,
		msgchan:  make(chan InnerMessage),
	}

	return dbmq
}

func (dbmq *DBMessageQueue) SyncSchema() error {

	var err error
	tables := []interface{}{
		new(po.MessageQueue),
	}
	err = dbmq.dbclient.SyncTables(tables)
	return err
}

func (dbmq *DBMessageQueue) SetLogger(logger *slog.Logger) {
	dbmq.logger = logger
}

// CleanExecutionMessage  清理ExecutionMessage
func (dbmq *DBMessageQueue) CleanExecutionMessage(execution_id int) error {

	var err error
	tx, maker := dbmq.dbclient.NewTxMaker(nil)
	defer maker.Close(&err)
	tx = tx.Where("execution_id = ?", execution_id).Delete(new(po.MessageQueue))
	err = tx.Error
	return err
}

// SetOption  设置option选项
func (dbmq *DBMessageQueue) SetOption(opt DBQueueOption) {
	dbmq.option = opt
}

// SendInnerMessage Send InnerMessage ,Save message  in DB
func (mq *DBMessageQueue) SendInnerMessage(message InnerMessageBody, sendTime *time.Time) error {

	var err error

	tx, maker := mq.dbclient.NewTxMaker(nil)
	defer maker.Close(&err)
	var dbMessage = po.MessageQueue{
		ExecutionID: message.ExecutionID,
		Type:        message.Type,
		StepID:      message.StepID,
		Data:        string(message.Data),
		Status:      DBMessageStatus.Created,
		SendTime:    time.Now(),
		Topic:       message.Class,
	}
	if sendTime != nil {
		dbMessage.SendTime = *sendTime
	}

	tx = tx.Create(&dbMessage)
	err = tx.Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	mq.logger.Debug(fmt.Sprintf("db messagequeue send message : %d ", dbMessage.ID))

	return nil
}

// ReceiveInnerMessage return  chan with type InnerMessage
func (dbmq *DBMessageQueue) ReceiveInnerMessage() (<-chan InnerMessage, error) {

	return dbmq.msgchan, nil
}

// GetUnreadMessageIDsOnce 获得所有未处理的消息列表
func (mq *DBMessageQueue) GetUnreadMessageIDsOnce() ([]int, error) {

	var err error
	tx, maker := mq.dbclient.NewTxMaker(nil)
	defer maker.Close(&err)

	var MsgIDs []int
	err = tx.Model(&po.MessageQueue{}).Select("id").Where(
		"status = ? and send_time <= ? ", DBMessageStatus.Created, time.Now()).Order("id asc").
		Limit(mq.option.PollingLimit).Find(&MsgIDs).Error

	if err != nil {
		tx.Rollback()
		return []int{}, err
	}
	tx.Commit()
	return MsgIDs, nil
}

// dispatchInnerMessage 投递消息 到通道
func (mq *DBMessageQueue) dispatchInnerMessage(msgID int, innerqueue InnerMessageQueue) error {
	var err error
	freshMsg := po.MessageQueue{}

	tx, maker := mq.dbclient.NewTxMaker(nil)
	defer maker.Close(&err)
	// 加排它锁， 处理activity task
	txf := rdb.ForUpdate(tx)
	err = txf.Take(&freshMsg, msgID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if freshMsg.Status != DBMessageStatus.Created {
		tx.Rollback()
		err = fmt.Errorf("db message status is invalid, id : %d, status: %s", msgID, freshMsg.Status)
		return err
	}

	// forward message to master queue
	newmsg := InnerMessageBody{
		ExecutionID: freshMsg.ExecutionID,
		StepID:      freshMsg.StepID,
		Class:       freshMsg.Topic,
		Type:        freshMsg.Type,
		Data:        freshMsg.Data,
	}

	err = innerqueue.SendInnerMessage(newmsg, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 处理成功， 更新状态

	tx = tx.Model(po.MessageQueue{ID: freshMsg.ID}).Delete(&po.MessageQueue{})
	err = tx.Error

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (mq *DBMessageQueue) StartPolling(masterqueue InnerMessageQueue) {
	mq.wg.Add(1)
	go mq.PollingInnerMessage(masterqueue)
}

// PollingInnerMessage polling inner message from db and send to channel
// masterqueue is the queue that will receive the message

func (mq *DBMessageQueue) PollingInnerMessage(masterqueue InnerMessageQueue) {
	mq.logger.Info("PollingInnerMessage Started.")

	var fetchmsg = func() {
		defer mq.wg.Done()

		msgids, err := mq.GetUnreadMessageIDsOnce()
		if err != nil {
			mq.logger.Error("get unread message ids failed", "error", err.Error())
			return
		}
		for _, msgID := range msgids {
			select {
			case <-mq.ctx.Done():
				return
			default:
				err = mq.dispatchInnerMessage(msgID, masterqueue)
				if err != nil {
					mq.logger.Error("dispatch innermessage failed", "error", err.Error())
				}
			}
		}
	}

	ticker := time.NewTicker(mq.option.PollingDuration)
	defer ticker.Stop()

	for {
		select {
		case <-mq.ctx.Done():
			return
		case <-ticker.C:
			mq.wg.Add(1)
			fetchmsg()
		}
	}
}

// Close stop async task message queue
func (mq *DBMessageQueue) Close() error {

	mq.cancel()
	mq.wg.Wait()
	return nil
}

// ParseDBQueueOption parse
func ParseDBQueueOption(data map[string]string) (DBQueueOption, error) {

	var configmapfuncs = map[string]func(*DBQueueOption, string) error{
		"pollingduration": func(c *DBQueueOption, v string) error {
			var err error
			c.PollingDuration, err = time.ParseDuration(v)
			return err
		},
		"limit": func(c *DBQueueOption, v string) error {
			var err error
			c.PollingLimit, err = strconv.Atoi(v)
			return err
		},
	}

	var err error
	var option = DefaultDBDelayQueueOption
	err = config.ParseMapStringConfig(&option, data, configmapfuncs)
	if err != nil {
		return DBQueueOption{}, err
	}
	return option, nil
}
