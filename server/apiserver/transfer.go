package apiserver

import (
	"time"

	"github.com/mmtbak/microlibrary/paging"
	"github.com/skyflow-workflow/skyflow_backbend/gen/pb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
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

func ToPBExecution(in po.Execution) *pb.ExecutionListItem {

	resp := &pb.ExecutionListItem{
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
		resp.FinishTime = in.StartTime.Unix()
	}
	return resp
}

func ToVOPageRequest(req *pb.PageRequest) paging.PageRequest {

	var voreq = paging.DefaultPageRequest
	if req != nil {
		if req.PageSize > 0 {
			voreq.PageSize = int(req.PageSize)
		}
		// 单页不能超过最大值
		if req.PageSize > int64(MaxPageSize) {
			voreq.PageNumber = MaxPageSize
		}
		if req.PageNumber > 0 {
			voreq.PageNumber = int(req.PageNumber)
		}
	}

	return voreq
}

func ToPBPageResponse(req paging.PageResponse) *pb.PageResponse {

	var resp = &pb.PageResponse{
		PageSize:   int64(req.PageSize),
		PageNumber: int64(req.PageNumber),
		Count:      int64(req.Count),
		PageCount:  int64(req.PageCount),
	}
	return resp
}

func ToPBExecutionEvent(in po.ExecutionEvent) *pb.ExecutionEventInfo {

	resp := &pb.ExecutionEventInfo{
		StateId:    int64(in.StateID),
		StateName:  in.StateName,
		EventType:  in.EventType,
		CreateTime: in.CreateTime.String(),
		StartTime:  in.StartTime.String(),
		FinishTime: ToTimeString(in.FinishTime),
		Data:       in.Data,
	}
	return resp

}

func ToPBNamespace(in po.Namespace) *pb.NamespaceListItem {
	resp := &pb.NamespaceListItem{
		Name:       in.Name,
		Comment:    in.Comment,
		CreateTime: in.CreateTime.Unix(),
		UpdateTime: in.UpdateTime.Unix(),
	}
	return resp
}

func ToPBActivity(in po.Activity) *pb.ActivityListItem {
	resp := &pb.ActivityListItem{
		Name:        in.Name,
		Comment:     in.Comment,
		ActivityUri: in.URI,
		CreateTime:  in.CreateTime.Unix(),
		UpdateTime:  in.UpdateTime.Unix(),
	}
	return resp
}
func ToPBStateMachine(in po.StateMachine) *pb.StateMachineListItem {
	resp := &pb.StateMachineListItem{
		Name:            in.Name,
		Comment:         in.Comment,
		StatemachineUri: in.URI,
		CreateTime:      in.CreateTime.Unix(),
		UpdateTime:      in.UpdateTime.Unix(),
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
	var out = []TO{}
	for _, d := range in {
		out = append(out, f(d))
	}
	return out
}
