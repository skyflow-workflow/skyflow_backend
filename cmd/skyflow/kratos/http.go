package kratos

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
	pbv1 "github.com/skyflow-workflow/skyflow_backend/api/v1"
	"github.com/skyflow-workflow/skyflow_backend/workflow/pberror"
)

// Response Http服务返回通用结构
type Response struct {
	Success      bool            `json:"success"`
	ErrorCode    string          `json:"error_code"`
	ReturnCode   int             `json:"return_code"`
	ErrorMessage string          `json:"error_message"`
	Data         json.RawMessage `json:"data"`
}

// CustomResponseEncoder 自定义响应编码器
func CustomResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	var err error
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 如果已经是错误响应，直接返回
	if err, ok := v.(error); ok {
		return err
	}

	// 包装成功响应
	var resp = Response{
		Success:      true,
		ErrorMessage: "",
		ErrorCode:    "",
		ReturnCode:   0,
	}
	switch data := v.(type) {
	case []byte:
		resp.Data = data
	default:
		vdata, err := json.Marshal(v)
		if err != nil {
			return err
		}
		resp.Data = vdata
	}
	// 序列化响应
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		slog.Error("write response error: %v", err)
		return err
	}
	return nil
}

// CustomErrorHandler 自定义错误处理器
func CustomErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var resp = Response{
		Success:      false,
		ErrorMessage: err.Error(),
		ErrorCode:    "",
		Data:         nil,
		ReturnCode:   int(-1),
	}
	var pberr *pberror.PBError
	var kerr *errors.Error
	var ok bool
	var statusCode int
	// kratos 内部错误
	kerr = errors.FromError(err)
	if kerr == nil {
		slog.Error("unknown error", "error", err)
		// 未知错误
		w.WriteHeader(http.StatusInternalServerError)
		resp.ErrorCode = pbv1.ErrorCode_UnknownError.String()
		resp.ReturnCode = int(pbv1.ErrorCode_UnknownError)
		resp.ErrorMessage = err.Error()
		goto WRITE_MESSAGE
	}
	pberr, ok = pberror.IsPBError(err)
	if ok {
		// 业务错误
		slog.Error("pberr error", "error", kerr)
		resp.ErrorCode = pberr.Code.String()
		resp.ErrorMessage = pberr.Message
		resp.ReturnCode = int(pberr.Code)
		w.WriteHeader(http.StatusOK)
		goto WRITE_MESSAGE
	}
	slog.Error("kratos error", "error", kerr)
	statusCode = getStatusCode(int(kerr.Code))
	resp.ErrorCode = kerr.Reason
	resp.ErrorMessage = kerr.Message
	w.WriteHeader(statusCode)
	goto WRITE_MESSAGE
WRITE_MESSAGE:
	werr := json.NewEncoder(w).Encode(resp)
	if werr != nil {
		slog.Error("write response error", "error", werr)
	}
}

// getStatusCode 根据业务错误码获取 HTTP 状态码
func getStatusCode(bizCode int) int {
	switch bizCode {
	case 400:
		return http.StatusBadRequest
	case 401:
		return http.StatusUnauthorized
	case 403:
		return http.StatusForbidden
	case 404:
		return http.StatusNotFound
	case 408:
		return http.StatusRequestTimeout
	case 500:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
