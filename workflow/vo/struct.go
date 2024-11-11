package vo

import "fmt"

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
