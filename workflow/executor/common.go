package executor

import "encoding/json"

// MergeExecutionStatus parallel/map 下有多个子流程，需要能合并多个子流程的状态成 parallel/map 的状态
func MergeExecutionStatus(status []string) _ExecutionStatusType {

	// 初始 = 最小的状态
	var finialstatus = ExecutionStatus.Created
	var finialweight = ExecutionStatusWeight[finialstatus]

	for _, st := range status {
		w, ok := ExecutionStatusWeight[_ExecutionStatusType(st)]
		if !ok {
			finialstatus = ExecutionStatus.Running
			return finialstatus
		}
		if w > finialweight {
			finialstatus = _ExecutionStatusType(st)
			finialweight = w
		}
		if finialstatus == ExecutionStatus.Running {
			break
		}
	}
	return finialstatus
}

// JSONString  serialize interface to string
func JSONString(v interface{}) (string, error) {
	var err error
	jsonbyte, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	jsonbytestr := string(jsonbyte)

	return jsonbytestr, nil
}
