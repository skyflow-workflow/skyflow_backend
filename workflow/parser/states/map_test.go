package states

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestParseMap(t *testing.T) {

	var testcases = []struct {
		template  string
		wantError bool
	}{
		{
			template: `
			{
				"Type": "Map",
				"InputPath": "$.detail",
				"ItemsPath": "$.shipped",
				"MaxConcurrency": 3,
				"ItemProcessor": {
				  "StartAt": "Validate",
				  "States": {
					"Validate": {
					  "Type": "Task",
					  "Resource": "arn:aws:lambda:us-east-1:123456789012:function:ship-val",
					  "End": true
					}
				  }
				},
				"ResultPath": "$.detail.shipped",
				"End": true
			  }	`,
			wantError: false,
		},
		{
			template: `
			{
				"ItemResultPath": "$.xyz",
				"MaxConcurrency": 1,
				"ItemProcessor": {
					"StartAt": "MHello",
					"States": {
						"MHello": {
							"Type": "Pass",
							"Result": "Hello",
							"Next": "MWorld"
						},
						"MWorld": {
							"End": true,
							"Type": "Pass",
							"Result": "World"
						}
					}
				},
				"Next": "ParallelState",
				"Type": "Map",
				"InputPath": "$.detail",
				"ItemsPath": "$.shipped"
			}`,
			wantError: false,
		},
		{
			template: `
			{
				"Type": "Map",
				"ItemsPath": "$.BNIpList",
				"MaxConcurrency": 0,
				"ItemProcessor": {
					"StartAt": "InstallOnePkgMWByOne",
					"States": {
						"InstallOnePkgMWByOne": {
							"Type": "Task",
							"Resource": "activity:yotta_dev/InstallOnePkgMWByOne",
							"Parameters": {
								"task_id.$": "$.task_id",
								"ip.$": "$.ip"
							},
							"OutputPath": "$",
							"ResultPath": "",
							"End": true,
							"Comment": "BlobNode安装部署包"
						}
					},
					"OutputPath": "$",
					"ResultPath": "",
					"Next": "BNEnableNodeMW",
					"Comment": "BlobNode安装部署包"
				}
			}
			`,
			wantError: true,
		},
	}
	for idx, tt := range testcases {
		fmt.Println("index :", idx)
		state, err := NewMapStateFromString(tt.template, StartDepth)
		fmt.Println(err)
		assert.Equal(t, tt.wantError, err != nil)
		if err == nil {
			state.SetName("Mapname")
			bone := state.GetBone()
			bonestr, _ := json.Marshal(bone)
			fmt.Println(string(bonestr))
		}
	}
}

func TestMapItem(t *testing.T) {

	var err error

	statejson := `
	{
		"Type": "Map",
		"InputPath": "$.detail",
		"ItemsPath": "$.shipped",
		"ItemResultPath" : "$.x",
		"MaxConcurrency": 0,
		"ItemProcessor": {
			"StartAt": "MHello",
			"States": {
				"MHello": {
					"Type": "Pass",
					"Result": "Hello",
					"Next": "MWorld"
				},
				"MWorld": {
					"Type": "Pass",
					"Result": "World",
					"End": true
				}
			}
		},
		"Next": "ParallelState"
	}
	`
	inputjson := `
	{
		"ship-date": "2016-03-14T01:59:00Z",
		"detail": {
			"delivery-partner":"UQS",
			"shipped": [{
				"prod": "R31",
				"dest-code": 9511,
				"quantity": 1344
			}, {
				"prod": "S39",
				"dest-code": 9511,
				"quantity": 40
			}, {
				"prod": "R31",
				"dest-code": 9833,
				"quantity": 12
			}, {
				"prod": "R40",
				"dest-code": 9860,
				"quantity": 887
			}, {
				"prod": "R40",
				"dest-code": 9860,
				"quantity": 888
			}, {
				"prod": "R40",
				"dest-code": 9860,
				"quantity": 889
			}, {
				"prod": "R40",
				"dest-code": 9511,
				"quantity": 1220
			}]
		}
	}
	`

	var mapstruct = map[string]interface{}{}
	var minput = map[string]interface{}{}
	json.Unmarshal([]byte(statejson), &mapstruct)
	json.Unmarshal([]byte(inputjson), &minput)
	mapstate, err := NewMapStateFromMap(mapstruct, 1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	items, err := mapstate.GetIteratorItems(minput)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(len(items))
	for _, item := range items {
		fmt.Println(item)
	}
	jsonbyte, err := json.Marshal(items)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(jsonbyte))

}

// /usr/bin/go test -benchmem -run=^$ -bench ^BenchmarkMapGetIteratorItems$ gopkg.mihoyo.com/plat/cloudflow/workflow/parser/grammar
func BenchmarkMapGetIteratorItems(b *testing.B) {

	var err error
	mapdefintion := `
	{
		"Name": "业务组模板变更阶段",
		"Type": "Map",
		"Next": "",
		"End": true,
		"Comment": "",
		"InputPath": "",
		"FailContiune": true,
		"ItemResultPath": "$.group_config",
		"ItemsPath": "$.biz_group_configs",
		"MaxConcurrency": 20,
		"ItemProcessor": {
		  "Comment": "Comment说明",
		  "Name": "业务组模板变更",
		  "StartAt": "查询发布分组ID",
		  "States": {
			"查询发布分组ID": {
			  "Comment": "通过参数查询发布分组ID信息",
			  "Parameters": {
				"app_name.$": "$.group_config.app_name",
				"biz_id.$": "$.biz_id",
				"cluster_tag.$": "$.cluster_tag",
				"env.$": "$.env",
				"service_group_name.$": "$.group_config.service_group_name"
			  },
			  "Resource": "activity:mockclouddeploy/QueryServiceGroup",
			  "ResultPath": "$.service_group",
			  "Type": "Task",
			  "Next": "变更应用分组环境配置"
			},
			"变更应用分组环境配置": {
			  "Name": "变更应用分组环境配置",
			  "OutputPath": "",
			  "Parameters": {
				"config.$": "$.group_config.config",
				"service_group_id.$": "$.service_group.id"
			  },
			  "Resource": "activity:mockclouddeploy/UpdateServiceGroupEnvironment",
			  "Type": "Task",
			  "End":true
			}
		  }
		}
	  }
	`

	fmt.Println("start benchmark")
	fmt.Println("map defintion len: ", len(mapdefintion))

	inputfile := "./map_input.json"
	input_content, err := os.ReadFile(inputfile)
	if err != nil {
		b.Error(err)
		return
	}
	fmt.Printf("input content size : %d\n", len(input_content))
	var inputobj = map[string]interface{}{}
	err = json.Unmarshal(input_content, &inputobj)
	if err != nil {
		b.Error(err)
		return
	}
	// fmt.Println("input object size : ", size.Of(inputobj))

	mapstate, err := NewMapStateFromString(mapdefintion, 1)
	if err != nil {
		b.Error(err)
		return
	}
	iterobjects, err := mapstate.GetIteratorItems2(inputobj)
	if err != nil {
		b.Error(err)
		return
	}
	// fmt.Printf("iterator objects size : %d\n", size.Of(iterobjects))
	iterobjectsbyte, err := json.Marshal(iterobjects)
	if err != nil {
		b.Error(err)
		return
	}
	fmt.Println("iterator objects string  len: ", len(iterobjectsbyte))
	// var iterobjectsstr string
	// iterobjectsstr := string(iterobjectsbyte)
	// fmt.Println(iterobjectsstr)

}
