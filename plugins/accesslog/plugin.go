//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 THL A29 Limited, a Tencent company.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the  Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

// Package accesslog contains proxy request logging.
package accesslog

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/trace"
	"trpc.group/trpc-go/trpc-gateway/common/convert"
	gerrs "trpc.group/trpc-go/trpc-gateway/common/errs"
	"trpc.group/trpc-go/trpc-gateway/common/gwmsg"
	"trpc.group/trpc-go/trpc-gateway/common/http"
	cplugin "trpc.group/trpc-go/trpc-gateway/common/plugin"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/filter"
	trpcHttp "trpc.group/trpc-go/trpc-go/http"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/naming/registry"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const pluginName = "accesslog"

func init() {
	plugin.Register(pluginName, log.DefaultLogFactory)
	plugin.Register(pluginName, &Plugin{})
}

// Plugin defines the plugin.
type Plugin struct {
}

// Type returns the plugin type.
func (p *Plugin) Type() string {
	return cplugin.DefaultType
}

// Setup initializes the plugin.
func (p *Plugin) Setup(string, plugin.Decoder) error {
	// Register the plugin
	filter.Register(pluginName, ServerFilter, nil)
	return nil
}

// Options represents the plugin configuration.
type Options struct {
	// TODO 如何保持有序？
	FieldList []map[string]string `yaml:"field_list"`
}

var DefaultOptions = &Options{
	FieldList: []map[string]string{},
}

// CheckConfig validates the plugin configuration and returns the parsed configuration object. Used in the ServerFilter
// method for parsing.
func (p *Plugin) CheckConfig(_ string, decoder plugin.Decoder) error {
	options := &Options{}
	if err := decoder.Decode(options); err != nil {
		return gerrs.Wrap(err, "decode access log config error")
	}
	DefaultOptions = options
	return nil
}

// ServerFilter represents the proxy information.
func ServerFilter(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (interface{}, error) {
	starttime := time.Now()
	rsp, err := next(ctx, req)
	finishtime := time.Now()
	d := finishtime.Sub(starttime)
	accessLog(ctx, err, d)
	return rsp, err
}

func accessLog(ctx context.Context, err error, d time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			log.ErrorContextf(ctx, "access_log_panic:%s,stack:%s", r, debug.Stack())
		}
	}()

	msg := trpc.Message(ctx)
	fieldList := []log.Field{
		// Interface path, corresponding to the "method" field configured in router.yaml
		{Key: "server_method", Value: msg.CallerMethod()},
		{Key: "service_name", Value: msg.CalleeServiceName()},
		{Key: "server", Value: msg.CalleeServer()},
		{Key: "server_rpc_name", Value: msg.ServerRPCName()},
		// Full interface path, not empty when the interface has rewriting
		{Key: "upstream_path", Value: msg.CallerMethod()},
		{Key: "err_no", Value: fmt.Sprint(errs.Code(err))},
		{Key: "err_msg", Value: getErrMSG(err)},
		{Key: "local_ip", Value: trpc.GlobalConfig().Global.LocalIP},
		// Backend service ip:port
		{Key: "remote_addr", Value: msg.RemoteAddr().String()},
		{Key: "upstream_service", Value: msg.CallerServiceName()},
		{Key: "upstream_protocol", Value: getUpstreamProtocol(ctx)},
		{Key: "elapsed", Value: d.String()},
		{Key: "trace_id", Value: getTraceID(ctx)},
		{Key: "request_id", Value: getRequestID(ctx)},
	}
	httpHeader := trpcHttp.Head(ctx)
	if httpHeader != nil {
		// HTTP request header
		httpFieldList := []log.Field{
			log.Field{Key: "protocol", Value: httpHeader.Request.Proto},
			log.Field{Key: "http_method", Value: httpHeader.Request.Method},
			log.Field{Key: "path", Value: string(httpHeader.Request.RequestURI)},
		}
		fieldList = append(fieldList, httpFieldList...)

	}

	// If it is an HTTP proxy, try to get the common parameters for reporting
	fctx := http.RequestContext(ctx)
	if fctx != nil {
		fasthttpFieldList := []log.Field{
			// Router ID
			{Key: "router_id", Value: gwmsg.GwMessage(ctx).RouterID()},
			{Key: "upstream_status", Value: fmt.Sprint(fctx.Response.StatusCode())},
			{Key: "remote_addr", Value: getClientIPFromContext(fctx)},
			{Key: "user_agent", Value: string(fctx.Request.Header.UserAgent())},
			{Key: "host", Value: string(fctx.Host())},
			{Key: "referer", Value: string(fctx.Referer())},
			{Key: "server_protocol", Value: string(fctx.Request.Header.Protocol())},
		}
		fieldList = append(fieldList, fasthttpFieldList...)
	}

	extFieldList, err := getExtFields(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "get ext accesslog field err:%s", err)
	}
	fieldList = append(fieldList, extFieldList...)
	// Get custom business fields
	// fieldList = append(fieldList, DefaultBusinessFields(ctx, node)...)
	getLogger().With(fieldList...).Info("accesslog")
}

// BusinessFields sets custom business log fields
type BusinessFields func(ctx context.Context, node *registry.Node) []log.Field

// DefaultBusinessFields overrides this method to set custom business log fields
var DefaultBusinessFields BusinessFields = func(ctx context.Context, _ *registry.Node) []log.Field {
	return nil
}

// Get extension fields
func getExtFields(ctx context.Context) ([]log.Field, error) {
	fctx := http.RequestContext(ctx)
	options := DefaultOptions
	log.DebugContextf(ctx, "accesslog_config:%s", convert.ToJSONStr(options))
	// Iterate through the configuration to get business parameters
	var fieldList []log.Field
	for _, fieldMap := range options.FieldList {
		for key, fieldName := range fieldMap {
			fieldList = append(fieldList,
				log.Field{
					Key:   key,
					Value: http.GetString(fctx, fieldName),
				},
			)
		}
	}
	return fieldList, nil
}

// Get the protocol
func getUpstreamProtocol(ctx context.Context) string {
	if gwmsg.GwMessage(ctx).TargetService() == nil {
		return ""
	}
	return gwmsg.GwMessage(ctx).TargetService().Protocol
}

func getRequestID(ctx context.Context) string {
	msg := codec.Message(ctx)
	if msg.RequestID() == 0 {
		return "" // If the request ID is not set, return an empty string
	}
	return fmt.Sprintf("%d", msg.RequestID())

}

// Get the full interface path
func getUpstreamPath(ctx context.Context) string {
	fctx := http.RequestContext(ctx)
	if fctx == nil {
		return ""
	}
	// Add reporting for the full path. For example: /a/{article_id}, msg.CallerMethod() is /a/, need to add reporting for
	// /a/{article_id} for troubleshooting purposes
	if string(fctx.Path()) != codec.Message(ctx).CallerMethod() {
		return string(fctx.Path())
	}
	return ""
}

// Get the error code
func getErrMSG(err error) string {
	if err == nil {
		return ""
	}
	return errs.Msg(err)
}

// Get the trace ID
func getTraceID(ctx context.Context) string {
	span := trace.SpanContextFromContext(ctx)
	if span.IsValid() {
		return span.TraceID().String()
	}
	return ""
}

// Get the logger
func getLogger() log.Logger {
	logger := log.GetDefaultLogger()
	if l := log.Get(pluginName); l != nil {
		logger = l
	}
	return logger
}

const localAddress = "127.0.0.1"

// GetClientIPFromContext retrieves the user's IP address from the context
func getClientIPFromContext(fCtx *fasthttp.RequestCtx) string {
	clientIPByte := fCtx.Request.Header.Peek(fasthttp.HeaderXForwardedFor)
	clientIPs := strings.Split(string(clientIPByte), ",")
	if len(clientIPs) == 0 {
		return ""
	}

	clientIP := clientIPs[0]
	if len(clientIP) > 0 && clientIP != localAddress {
		return clientIP
	}
	clientIP = strings.TrimSpace(string(fCtx.Request.Header.Peek("X-Real-Ip")))
	if clientIP != "" && clientIP != localAddress {
		return clientIP
	}
	return ""
}
