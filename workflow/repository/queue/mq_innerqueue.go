package queue

import (
	"context"
	"log/slog"
	"time"

	"github.com/mmtbak/microlibrary/mq"
)

// MQInnerMessageQueue mq based inner message queue
type MQInnerMessageQueue struct {
	basequeue mq.MessageQueue
	msgchan   chan InnerMessage
	ctx       context.Context
	cancel    context.CancelFunc
	logger    *slog.Logger
}

func NewMQInnerMessageQueue(queue mq.MessageQueue) *MQInnerMessageQueue {

	ctx, cancel := context.WithCancel(context.Background())

	innerMQ := &MQInnerMessageQueue{
		basequeue: queue,
		msgchan:   make(chan InnerMessage),
		cancel:    cancel,
		ctx:       ctx,
		logger:    slog.Default(),
	}
	return innerMQ

}

func (imq *MQInnerMessageQueue) SendInnerMessage(msgbody InnerMessageBody, sendtime *time.Time) error {

	var err error
	sendOpt := mq.NewSendMsgOption()
	if sendtime != nil {
		sendOpt = sendOpt.WithSendtime(*sendtime)
	}
	data, err := msgbody.Marshal()
	if err != nil {
		return err
	}
	err = imq.basequeue.SendMessage(data, sendOpt)
	return err
}

// ReceiveInnerMessage  receive  inner message
func (imq *MQInnerMessageQueue) ReceiveInnerMessage() (<-chan InnerMessage, error) {

	basechan, err := imq.basequeue.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	go func() {
		for {

			select {
			case <-imq.ctx.Done():
				return
			case basemsg := <-basechan:
				var err error
				var innermsgbody = &InnerMessageBody{}
				bodyData := basemsg.Body()
				innermsgbody, err = UnmarshalInnerMessageBody(bodyData)
				if err != nil {
					imq.logger.Error("unmarshal inner message body error", "error", err)
					continue
				}
				innermsg := &MQInnerMessage{
					msg:  basemsg,
					body: innermsgbody,
					mq:   imq,
				}
				imq.msgchan <- innermsg
			}
		}
	}()

	return imq.msgchan, nil
}

func (imq *MQInnerMessageQueue) SyncSchema() error {
	err := imq.basequeue.SyncSchema()
	return err
}

func (imq *MQInnerMessageQueue) CleanExecutionMessage(executionID int) error {
	return nil
}

func (imq *MQInnerMessageQueue) Close() error {
	imq.cancel()
	return imq.basequeue.Close()
}

type MQInnerMessage struct {
	mq   *MQInnerMessageQueue
	body *InnerMessageBody
	msg  mq.Message
}

// Body return body
func (im *MQInnerMessage) Body() InnerMessageBody {
	return *im.body
}

// Ack return nil
func (im *MQInnerMessage) Ack() error {
	return im.msg.Ack()
}

// Ack return nil
func (im *MQInnerMessage) ID() string {
	return im.msg.ID()
}

// Nack  return nil
func (im *MQInnerMessage) Nack() error {
	return im.msg.Nack()
}
