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

// Wait ...
type Wait struct {
	*BaseState
	*WaitBody
}

func (w *Wait) Init() error {
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
func (w *Wait) GetWakeupTime(input any) (time.Time, error) {

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
func (w *Wait) GetNextState(input interface{}) (NextState, error) {
	ns := NextState{
		Name:   w.Next,
		Output: input,
	}
	return ns, nil
}
