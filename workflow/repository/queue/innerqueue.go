package queue

import (
	"fmt"

	"github.com/mmtbak/microlibrary/config"
	"github.com/mmtbak/microlibrary/mq"
	"github.com/mmtbak/microlibrary/rdb"
)

// NewInnerMessageQueueFromConfig create new inner message queue
func NewInnerMessageQueueFromConfig(conf config.AccessPoint) (InnerMessageQueue, error) {

	dsnData, err := conf.Decode(nil)
	if err != nil {
		return nil, err
	}
	switch dsnData.Scheme {
	case "kafka":
		basequeue, err := mq.NewKafkaMessageQueue(conf)
		if err != nil {
			return nil, err
		}
		innermq := NewMQInnerMessageQueue(basequeue)

		return innermq, nil
	case "mysql":
		dbopt, err := ParseDBQueueOption(dsnData.Params)
		if err != nil {
			return nil, err
		}
		dbconfig, err := rdb.ParseConfig(conf)
		if err != nil {
			return nil, err
		}
		dbclient, err := rdb.NewDBClient(dbconfig)
		if err != nil {
			return nil, err
		}

		innermq := NewDBMessageQueue(dbclient, dbopt)
		return innermq, err
	default:
		err = fmt.Errorf("message queue :unsupported message queue  schema '%s'", dsnData.Scheme)
	}
	return nil, err
}
