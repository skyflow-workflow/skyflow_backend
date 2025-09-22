package states

import (
	"fmt"
	"time"

	"github.com/skyflow-workflow/skyflow_backbend/pkg/jsonpath"
)

var TimeFormatCommon = "2006-01-02 15:04:05"
var TimeFormatList = []string{TimeFormatCommon, time.RFC3339}

// WaitBody ...
type WaitBody struct {
	Seconds       uint   `mapstructure:"Seconds" validate:"gte=0"`
	Timestamp     string `mapstructure:"Timestamp" validate:"gte=0"`
	SecondsPath   string `mapstructure:"SecondsPath" validate:"gte=0"`
	TimestampPath string `mapstructure:"TimestampPath" validate:"gte=0"`
}

// WaitState ...
type WaitState struct {
	*BaseState
	*WaitBody
}

// NewWaitStateFromString Create New Wait State
func NewWaitStateFromString(definition string) (state *WaitState, err error) {

	data, err := ToMap(definition)
	if err != nil {
		return
	}
	state, err = NewWaitStateFromMap(data)
	return
}

// NewWaitStateFromMap Create New Wait State
func NewWaitStateFromMap(data map[string]interface{}) (state *WaitState, err error) {

	bs, err := NewBaseStateFromMap(data)
	if err != nil {
		return
	}
	body, err := NewWaitBodyFromMap(data)
	if err != nil {
		return
	}
	state = &WaitState{
		BaseState: bs,
		WaitBody:  body,
	}
	return
}

func NewWaitBodyFromMap(data map[string]interface{}) (state *WaitBody, err error) {
	state = &WaitBody{
		Seconds:       0,
		Timestamp:     "",
		SecondsPath:   "",
		TimestampPath: "",
	}

	err = MapStructDecode(data, state)
	if err != nil {
		return
	}

	return
}

// InitByMap Inititalize Wait Content
func InitWaitBodyByMap(body *WaitBody, data map[string]interface{}) (err error) {

	// 初始化自身
	err = MapStructDecode(data, body)
	if err != nil {
		return err
	}
	err = myValidate.Struct(body)
	if err != nil {
		return err
	}

	err = body.Init()
	if err != nil {
		return err
	}
	return nil
}

func (w *WaitState) GetBaseState() *BaseState {
	return w.BaseState
}

func (w *WaitBody) Init() error {
	var err error
	err = myValidate.Struct(w)
	if err != nil {
		return err
	}
	if w.Timestamp != "" {
		for _, format := range TimeFormatList {
			_, err = time.Parse(format, w.Timestamp)
			if err != nil {
				continue
			}
		}
		if err != nil {
			err = fmt.Errorf("field 'Timestamp' format error")
			return err
		}
	}
	if w.SecondsPath != "" {
		_, err = jsonpath.JsonPathCompile(w.SecondsPath)
		if err != nil {
			return fmt.Errorf("field 'SecondsPath' error : %w", err)
		}
	}
	if w.TimestampPath != "" {
		_, err = jsonpath.JsonPathCompile(w.TimestampPath)
		if err != nil {
			return fmt.Errorf("field 'TimestampPath' error : %w", err)
		}
	}
	return nil
}

// GetWakeupTime Get Wait State Wake up time
func (w *WaitBody) GetWakeupTime(input any) (time.Time, error) {

	var err error
	now := time.Now()
	var dest time.Time
	if w.Seconds > 0 {
		dest = now.Add(time.Duration(w.Seconds) * time.Second)
		return dest, nil
	}
	if w.SecondsPath != "" {
		intervalSeconds, err := jsonpath.JsonPathGetValue(w.SecondsPath, input)
		if err != nil {
			return dest, err
		}

		seconds, ok := intervalSeconds.(int)
		if !ok {
			err = fmt.Errorf("SecondPath value should be type int: [ %s ]", w.SecondsPath)
			return dest, err
		}
		if seconds > 0 {
			dest = now.Add(time.Duration(seconds) * time.Second)
		} else {
			dest = now
		}
		return dest, nil
	}

	if w.Timestamp != "" {
		for _, format := range TimeFormatList {
			dest, err = time.Parse(format, w.Timestamp)
			if err != nil {
				continue
			}
			break
		}
		if err != nil {
			err = fmt.Errorf("field 'Timestamp' format error")
			return now, err
		}
		return dest, nil
	}
	if w.TimestampPath != "" {
		timestamp, err := jsonpath.JsonPathGetValue(w.TimestampPath, input)
		if err != nil {
			return dest, err
		}
		timestampStr, ok := timestamp.(string)
		if !ok {
			err = fmt.Errorf("TimestampPath value should be type string: [ %s ]", w.TimestampPath)
			return dest, err
		}
		for _, format := range TimeFormatList {
			dest, err = time.Parse(format, timestampStr)
			if err != nil {
				continue
			}
			break
		}
		if err != nil {
			err = fmt.Errorf("TimestampPath format error")
			return now, err
		}
		return dest, nil
	}

	return dest, nil
}

// GetNextState Get Next State
func (w *WaitState) GetNextState(input interface{}) (NextState, error) {
	ns := NextState{
		Name:   w.Next,
		Output: input,
	}
	return ns, nil
}
