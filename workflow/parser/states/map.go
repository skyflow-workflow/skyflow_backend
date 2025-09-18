package states

/*
 * @Author: mumangtao@gmail.com
 * @Date: 2019-12-06 10:58:30
 * @Last Modified by: mumangtao@gmail.com
 * @Last Modified time: 2020-08-09 01:31:43
 */

import (
	"encoding/json"
	"fmt"

	"github.com/mohae/deepcopy"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

//  Map Sate Parser Here
//  State Define
//"Validate-All": {
// 	"Type": "Map",
// 	"InputPath": "$.detail",
// 	"ItemsPath": "$.shipped",
// 	"MaxConcurrency": 0,
// 	"ItemProcessor": {
// 	  "StartAt": "Validate",
// 	  "States": {
// 		"Validate": {
// 		  "Type": "Task",
// 	      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:ship-val",
// 		  "End": true
// 		}
// 	  }
// 	},
// 	"ResultPath": "$.detail.shipped",
// 	"End": true
//   }
//  input
// {
// 	"ship-date": "2016-03-14T01:59:00Z",
// 	"detail": {
// 	  "delivery-partner": "UQS",
// 	  "shipped": [
// 		{ "prod": "R31", "dest-code": 9511, "quantity": 1344 },
// 		{ "prod": "S39", "dest-code": 9511, "quantity": 40 },
// 		{ "prod": "R31", "dest-code": 9833, "quantity": 12 },
// 		{ "prod": "R40", "dest-code": 9860, "quantity": 887 },
// 		{ "prod": "R40", "dest-code": 9511, "quantity": 1220 }
// 	  ]
// 	}
// }

// MapState Map State Struct
// 新增字段ItemResultPath，用来将Item 放置到Input中去的路径
type MapState struct {
	*BaseState `mapstructure:",squash"`
	*MapBody   `mapstructure:",squash"`
}

type MapBody struct {
	ItemsPath      string `mapstructure:"ItemsPath" validate:"required,gt=0,startswith=$"`
	ItemResultPath string `mapstructure:"ItemResultPath" validate:"gte=0"`
	MaxConcurrency int    `mapstructure:"MaxConcurrency" validate:"gte=0"`
	//FailContinue 一条分支执行结束后，是否自动启动新分支
	FailContinue  bool                   `mapstructure:"FailContinue" `
	ItemProcessor map[string]interface{} `mapstructure:"ItemProcessor" validate:"required,gt=0"`
}

// NewMapStateFromString New Map Struct
func NewMapStateFromString(definition string, depth int) (state *MapState, err error) {

	state, err = NewGenericGroupStateFromString(definition, depth, NewMapStateFromMap)
	return state, err
}

// NewMapStateFromMap New Map Struct
func NewMapStateFromMap(data map[string]interface{}, depth int) (*MapState, error) {

	var err error
	// 默认map初始化结构
	var state = &MapState{
		// BaseState:      NewDefautBaseState(),
		ItemsPath:      "$",
		ItemResultPath: "",
		MaxConcurrency: 100,
		FailContinue:   false,
		_Depth:         depth,
	}

	// basestate
	bs, err := NewBaseStateFromMap(data)
	if err != nil {
		return nil, err
	}
	state.BaseState = bs

	err = state.InitByMap(data)
	if err != nil {
		return nil, err
	}
	return state, err
}

// InitByMap Initialize Parser Map State Content
func (s *MapState) InitByMap(data map[string]interface{}) error {
	var err error

	// 初始化自身
	err = MapStructDecode(data, s)
	if err != nil {
		return err
	}
	err = myvalidate.Struct(s)
	if err != nil {
		return err
	}

	subsg, err := NewStateGroupFromMap(s.ItemProcessor, s._Depth+1)
	if err != nil {
		return err
	}

	s._ItemProcessor = subsg

	err = s.Init()
	if err != nil {
		return err
	}
	return nil
}

// 其他地方可以自己组装， 自己做初始化
func (s *MapState) Init() error {
	var err error
	if s.ItemsPath != "" {
		s._ItemsPath, err = jsonpath.Compile(s.ItemsPath)
		if err != nil {
			return fmt.Errorf("field 'ItemsPath' error : %w", err)
		}
	}

	// MaxConcurrency 限制
	// 只能放在这里， 因为有map 和 jobmap两种类型的初始化， 都会走init()
	if s.MaxConcurrency > quota.MaxMapConncurrencyLimit {
		return vo.ErrorMapConncurrencyLimitExceeded
	}

	// 初始化ItemProcessor
	smbone := s._ItemProcessor.GetBone()
	s._Bone.Branches = []StateBone{smbone}
	return nil
}

// GetIteratorItems Gen input for branch
func (s *MapState) GetIteratorItems(input interface{}) ([]interface{}, error) {

	var err error
	var ok bool
	var newoutput interface{}
	var items []interface{}

	var newinput interface{}
	newinput, err = s.GetInput(input)
	if err != nil {
		return items, err
	}

	// itempath 不为空
	if s.ItemsPath != "" {
		newoutput, err = jsonpath.JsonPathLookup(newinput, s.ItemsPath)
		if err != nil {
			return items, err
		}
		// items 必须是数组
		items, ok = newoutput.([]interface{})
		if !ok {
			err = fmt.Errorf("input iterator items  Should Array")
			return items, err
		}
	}
	// Item在Input中的路径, 如果不为空，则将Item放置到 ItemResultPath 路径中去。
	if s.ItemResultPath != "" {
		var newitems = []interface{}{}
		for idx := range items {
			copyinput := deepcopy.Copy(newinput)
			// copyitem := deepcopy.Copy(items[idx])
			// newitem, err := jsonpath.ReferencePathSetValue(copyinput, s.ItemResultPath, copyitem)
			newitem, err := jsonpath.ReferencePathSetValue(copyinput, s.ItemResultPath, items[idx])
			if err != nil {
				return items, err
			}
			newitems = append(newitems, newitem)
		}
		items = newitems
	}
	return items, nil
}

// GetIteratorItems2 Gen input for branch
// GetIteratorItems2 替代 GetIteratorItems的优化版本，减少了数据复制的次数
func (s *MapState) GetIteratorItems2(input interface{}) ([]string, error) {

	var err error
	var ok bool
	var newoutput interface{}
	var itemsresult []string
	var items []interface{}

	var newinput interface{}
	newinput, err = s.GetInput(input)
	if err != nil {
		return itemsresult, err
	}

	// itempath 不为空
	if s.ItemsPath != "" {
		newoutput, err = jsonpath.JsonPathLookup(newinput, s.ItemsPath)
		if err != nil {
			return itemsresult, err
		}
		// items 必须是数组
		items, ok = newoutput.([]interface{})
		if !ok {
			err = fmt.Errorf("input iterator items  Should Array")
			return itemsresult, err
		}
	}
	// Item在Input中的路径, 如果不为空，则将Item放置到 ItemResultPath 路径中去。
	if s.ItemResultPath != "" {
		// var itemsresult = []string{}
		for idx := range items {
			newitem, err := jsonpath.ReferencePathSetValue(newinput, s.ItemResultPath, items[idx])
			if err != nil {
				return itemsresult, err
			}
			newitembyte, err := json.Marshal(newitem)
			if err != nil {
				return itemsresult, err
			}
			itemsresult = append(itemsresult, string(newitembyte))
		}
		// items = newitems
	}
	return itemsresult, nil
}

// GetItemProcessor  获得需要Map流程模板
func (s *MapState) GetItemProcessor() *StateGroup {
	return s._ItemProcessor
}

func (s *MapState) SetDepth(depth int) {
	s._Depth = depth
}

func (s *MapState) SetName(name string) {
	s.BaseState.SetName(name)

}

func (s *MapState) SetItemProcessor(sg *StateGroup) {
	s._ItemProcessor = sg
}
