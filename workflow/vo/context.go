package vo

import "context"

var ContextKey = struct {
	RequestInfo string
}{
	RequestInfo: "RequestInfo",
}

func WithRequestInfo(ctx context.Context, reqinfo RequestInfo) context.Context {

	return context.WithValue(ctx, ContextKey.RequestInfo, reqinfo)
}

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
