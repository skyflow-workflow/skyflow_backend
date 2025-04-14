package vo

import (
	"context"
	"fmt"
)

// ContextKey ...
var ContextKey = struct {
	RequestInfo string
}{
	RequestInfo: "RequestInfo",
}

// WithRequestInfo ...
func WithRequestInfo(ctx context.Context, reqinfo RequestInfo) context.Context {

	return context.WithValue(ctx, ContextKey.RequestInfo, reqinfo)
}

// GetRequestInfo ...
func GetRequestInfo(ctx context.Context) RequestInfo {
	val := ctx.Value(ContextKey.RequestInfo)
	if val != nil {
		return val.(RequestInfo)
	}
	return RequestInfo{
		RemoteAddress: "",
		RequestType:   "unknown",
	}
}

// RequestInfo  发起请求的Request相关信息
type RequestInfo struct {
	// RemoteAddress 调用API的远程地址, ip:port
	RemoteAddress string
	// RequestType 请求方式： HTTP， GRPC
	RequestType string
}

// String format string
func (req RequestInfo) String() string {
	var res string
	res = fmt.Sprintf("RemoteAdress: %s, RequestType: %s ", req.RemoteAddress, req.RequestType)
	return res
}
