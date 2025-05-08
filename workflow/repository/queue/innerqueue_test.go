package queue

import (
	"fmt"
	"testing"

	"github.com/mmtbak/microlibrary/config"
)

func TestCreateInnerQueue(t *testing.T) {

	conf := config.AccessPoint{
		Source: "",
	}
	mq, err := NewInnerMessageQueueFromConfig(conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mq)
}
