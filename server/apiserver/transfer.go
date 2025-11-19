package apiserver

import (
	"time"

	"github.com/mmtbak/microlibrary/paging"
	pbv1 "github.com/skyflow-workflow/skyflow_backend/api/v1"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
)

const (
	// 时间格式
	timeformat = time.RFC3339
	// 分页单页 最大页大小
	MaxPageSize = 2000
)

func ToTimeString(t time.Time) string {
	return t.Format(timeformat)
}

func ToPBExecutionItem(in po.Execution) *pbv1.ExecutionListItem {

	resp := &pbv1.ExecutionListItem{
		ExecutionUuid: in.UUID,
		Status:        in.Status,
		Title:         in.Title,
		Definition:    in.Definition,
		CreateTime:    in.CreateTime.Unix(),
	}
	if in.StartTime != nil {
		resp.StartTime = in.StartTime.Unix()
	}
	if in.FinishTime != nil {
		resp.FinishTime = in.FinishTime.Unix()
	}
	return resp
}

func ToVOPageRequest(req *pbv1.PageRequest) paging.PageRequest {

	var voReq = paging.DefaultPageRequest
	if req != nil {
		if req.PageSize > 0 {
			voReq.PageSize = int(req.PageSize)
		}
		// 单页不能超过最大值
		if req.PageSize > int64(MaxPageSize) {
			voReq.PageSize = MaxPageSize
		}
		// 验证页码和页大小的合理性
		if req.PageNumber < 1 {
			voReq.PageNumber = 1
		}
		if req.PageSize < 1 {
			voReq.PageSize = 10 // 默认页大小
		}
		if req.PageNumber > 0 {
			voReq.PageNumber = int(req.PageNumber)
		}
	}

	return voReq
}

func ToPBPageResponse(req paging.PageResponse) *pbv1.PageResponse {

	var resp = &pbv1.PageResponse{
		PageSize:   int64(req.PageSize),
		PageNumber: int64(req.PageNumber),
		Count:      int64(req.Count),
		PageCount:  int64(req.PageCount),
	}
	return resp
}

func ToPBExecutionEvent(in po.ExecutionEvent) *pbv1.ExecutionEventInfo {

	resp := &pbv1.ExecutionEventInfo{
		StepId:     int64(in.StepID),
		StepName:   in.StepName,
		EventType:  in.EventType,
		CreateTime: in.CreateTime.String(),
		StartTime:  in.StartTime.String(),
		FinishTime: ToTimeString(in.FinishTime),
		Data:       in.Data,
	}
	return resp

}

func ToPBNamespace(in po.Namespace) *pbv1.NamespaceListItem {
	resp := &pbv1.NamespaceListItem{
		Name:        in.Name,
		Description: in.Description,
		CreateTime:  in.CreateTime.Unix(),
		UpdateTime:  in.UpdateTime.Unix(),
	}
	return resp
}

func ToPBActivityItem(in po.Activity) *pbv1.ActivityListItem {
	resp := &pbv1.ActivityListItem{
		Name:        in.Name,
		Description: in.Description,
		ActivityUri: in.URI,
		CreateTime:  in.CreateTime.Unix(),
		UpdateTime:  in.UpdateTime.Unix(),
	}
	return resp
}
func ToPBStateMachineItem(in po.StateMachine) *pbv1.StateMachineListItem {
	resp := &pbv1.StateMachineListItem{
		Name:            in.Name,
		Description:     in.Description,
		StatemachineUri: in.URI,
		CreateTime:      in.CreateTime.Unix(),
		UpdateTime:      in.UpdateTime.Unix(),
	}
	return resp
}
func ToPBStateMachine(in *po.StateMachine) *pbv1.StateMachineInfo {
	resp := &pbv1.StateMachineInfo{
		Name:            in.Name,
		Description:     in.Description,
		StatemachineUri: in.URI,
		CreateTime:      in.CreateTime.Unix(),
		UpdateTime:      in.UpdateTime.Unix(),
		Definition:      in.Definition,
	}
	return resp
}

// DataTransferArray transfer data from one array to another array
// in: input array
// f: transfer function
// out: output array
// example:
// in := []int{1, 2, 3}
//
//	f := func(i int) string {
//	    return strconv.Itoa(i)
//	}
//
// out := DataTransferArray(in, f)
// fmt.Println(out) // ["1", "2", "3"]
// out: []string
//
//	out := DataTransferArray(in, func(i int) string {
//	    return strconv.Itoa(i)
//	})
//
// fmt.Println(out) // ["1", "2", "3"]
func DataTransferArray[TI any, TO any](in []TI, f func(TI) TO) []TO {
	if len(in) == 0 {
		return []TO{}
	}
	// 预分配容量避免多次内存重分配
	out := make([]TO, 0, len(in))
	for _, d := range in {
		out = append(out, f(d))
	}
	return out
}
