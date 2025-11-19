package cmd

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestLoadConfig(t *testing.T) {

	conf, err := LoadConfig("./mock/skyflow.yaml")
	assert.Equal(t, err, nil)
	assert.Equal(t, conf.API.QPSLimit, 1000)

	svc, err := LoadService(conf)
	assert.Equal(t, err, nil)
	fmt.Println(svc)
}
