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
	dbClient     *rdb.DBClient
	wg           sync.WaitGroup
	option       DBQueueOption
	logger       *slog.Logger
	ctx          context.Context
	cancel       context.CancelFunc
	msgChan      chan InnerMessage
	forwardQueue InnerMessageQueue
	once         sync.Once
}

// NewDBMessageQueue Create New Inner DBMessage
func NewDBMessageQueue(client *rdb.DBClient, option DBQueueOption) *DBMessageQueue {

	ctx, cancel := context.WithCancel(context.Background())

	dbmq := &DBMessageQueue{
		dbClient: client,
		wg:       sync.WaitGroup{},
		ctx:      ctx,
		cancel:   cancel,
		option:   option,
		msgChan:  make(chan InnerMessage),
		logger:   slog.Default(),
		once:     sync.Once{},
	}

	return dbmq
}

func (dbmq *DBMessageQueue) SyncSchema() error {

	var err error
	tables := []interface{}{
		new(po.MessageQueue),
	}
	err = dbmq.dbClient.SyncTables(tables)
	return err
}

func (dbmq *DBMessageQueue) SetLogger(logger *slog.Logger) {
	dbmq.logger = logger
}

// CleanExecutionMessage  清理ExecutionMessage
func (dbmq *DBMessageQueue) CleanExecutionMessage(execution_id int) error {

	var err error
	tx, maker := dbmq.dbClient.NewTxMaker(nil)
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

	tx, maker := mq.dbClient.NewTxMaker(nil)
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

	return nil
}

// ReceiveInnerMessage return  chan with type InnerMessage
func (dbmq *DBMessageQueue) ReceiveInnerMessage() (<-chan InnerMessage, error) {

	return dbmq.msgChan, nil
}

// GetUnreadMessageIDsOnce 获得所有未处理的消息列表
func (mq *DBMessageQueue) GetUnreadMessageIDsOnce() ([]int, error) {

	var err error
	tx, maker := mq.dbClient.NewTxMaker(nil)
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

// DispatchInnerMessage dispatch message to forward queue
func (mq *DBMessageQueue) DispatchInnerMessage(msgID int) error {
	var err error
	freshMsg := po.MessageQueue{}

	tx := mq.dbClient.DB()
	// user non-blocking transaction
	conditionMessage := po.MessageQueue{
		ID:     int64(msgID),
		Status: DBMessageStatus.Created,
	}
	updateMessage := po.MessageQueue{
		Status: DBMessageStatus.Processed,
	}
	utx := tx.Where(conditionMessage).Updates(&updateMessage)
	if utx.Error != nil {
		return tx.Error
	}
	if utx.RowsAffected == 0 {
		err = fmt.Errorf("db message not found, id : %d", msgID)
		return err
	}
	err = tx.Model(&po.MessageQueue{}).Where("id = ?", msgID).First(&freshMsg).Error
	if err != nil {
		return fmt.Errorf("get message by id %d failed, error: %s", msgID, err.Error())
	}

	// forward message to forward queue
	newMsg := InnerMessageBody{
		ExecutionID: freshMsg.ExecutionID,
		StepID:      freshMsg.StepID,
		Class:       freshMsg.Topic,
		Type:        freshMsg.Type,
		Data:        freshMsg.Data,
	}

	if mq.forwardQueue != nil {
		err = mq.forwardQueue.SendInnerMessage(newMsg, nil)
		if err != nil {
			return err
		}
	} else {
		// if forwardQueue is nil, send message to msgChan
		// this is for testing purpose, in production, forwardQueue should not be nil
		mq.msgChan <- &DBQueueMessage{
			id:   int(freshMsg.ID),
			body: &newMsg,
		}
	}

	tx, maker := mq.dbClient.NewTxMaker(nil)
	defer maker.Close(&err)
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

func (mq *DBMessageQueue) SetForwardQueue(forwardQueue InnerMessageQueue) {
	mq.forwardQueue = forwardQueue
}

func (mq *DBMessageQueue) StartPolling() {
	go mq.once.Do(mq.RunPollingInnerMessage)
}

// PollingInnerMessage polling inner message from db and send to channel
// masterQueue is the queue that will receive the message

func (mq *DBMessageQueue) RunPollingInnerMessage() {
	mq.logger.Info("PollingInnerMessage Started.")

	var fetchmsg = func() {

		msgids, err := mq.GetUnreadMessageIDsOnce()
		if err != nil {
			mq.logger.Error("get unread message ids failed", "error", err.Error())
			return
		}
		for _, msgID := range msgids {
			err = mq.DispatchInnerMessage(msgID)
			if err != nil {
				mq.logger.Error("dispatch innermessage failed", "error", err.Error())
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
			fetchmsg()
		}
	}
}

// Close stop async task message queue
func (mq *DBMessageQueue) Close() error {
	close(mq.msgChan)
	mq.cancel()
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

type DBQueueMessage struct {
	id   int
	body *InnerMessageBody
}

func (t *DBQueueMessage) ID() string {
	return strconv.Itoa(t.id)
}
func (t *DBQueueMessage) Body() InnerMessageBody {
	return *t.body
}

func (t *DBQueueMessage) Ack() error {
	// DBQueueMessage does not need to ack, because it is already processed in DispatchInnerMessage
	return nil
}
func (t *DBQueueMessage) Nack() error {
	// DBQueueMessage does not need to nack, because it is already processed in DispatchInnerMessage
	return nil
}
