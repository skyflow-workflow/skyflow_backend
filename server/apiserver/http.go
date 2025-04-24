package apiserver

import (
	"bytes"
	"encoding/json"
	stdhttp "net/http"

	"trpc.group/trpc-go/trpc-go/errs"
	thttp "trpc.group/trpc-go/trpc-go/http"
)

func init() {
	thttp.DefaultServerCodec.ErrHandler = DefaultHTTPErrorHandler
	thttp.DefaultServerCodec.RspHandler = DefaultHttpRespHandler
}

// Response Http服务返回通用结构
type Response struct {
	Success    bool            `json:"success"`
	ErrorCode  string          `json:"error_code"`
	ReturnCode int             `json:"return_code"`
	ErrorMsg   string          `json:"error_message"`
	Data       json.RawMessage `json:"data"`
}

// DefaultHttpRespHandler  default trpc http handler
var DefaultHttpRespHandler = func(w stdhttp.ResponseWriter, r *stdhttp.Request, rspbody []byte) (err error) {
	var data json.RawMessage
	if len(rspbody) == 0 {
		data = []byte("null")
	} else {
		data = rspbody
	}
	bs, _ := json.Marshal(&Response{Success: true, ErrorCode: "", ErrorMsg: "", Data: data})
	_, err = bytes.NewBuffer(bs).WriteTo(w)
	return
}

// DefaultHTTPErrorHandler default http error handler
var DefaultHTTPErrorHandler = func(w stdhttp.ResponseWriter, r *stdhttp.Request, e *errs.Error) {
	// 填充指定格式错误信息到HTTP Body
	resp := Response{
		Success:    false,
		ErrorMsg:   e.Msg,
		ErrorCode:  "InternalError",
		Data:       nil,
		ReturnCode: int(e.Code),
	}
	bs, _ := json.Marshal(&resp)
	_, _ = bytes.NewBuffer(bs).WriteTo(w)
}
