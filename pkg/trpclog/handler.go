package trpclog

import (
	"context"
	"log/slog"

	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/log"
)

type Logger struct{}

type TrpcLogHandler struct {
	logger log.Logger
	groups []string
}

func ContextWithLogger(ctx context.Context, logger log.Logger) context.Context {
	msg := codec.Message(ctx)
	msg.WithLogger(logger)
	ctx = context.WithValue(ctx, codec.ContextKeyMessage, logger)
	return ctx
}

func GetLoggerFromContext(ctx context.Context) log.Logger {
	if l, ok := codec.Message(ctx).Logger().(log.Logger); ok {
		return l
	}
	return log.DefaultLogger
}

// HandlerFromTrpcLogger Create a slog.Handler from a trpc Logger.
func NewHandlerFromTrpcLogger(l log.Logger) *TrpcLogHandler {
	if l == nil {
		l = log.DefaultLogger
	}

	if ol, ok := l.(log.OptionLogger); ok {
		l = ol.WithOptions(log.WithAdditionalCallerSkip(2))
	}

	// 调整调用栈， 找到最原始打印日志的位置
	// l = l.CallSkip(3)
	h := &TrpcLogHandler{
		logger: l,
		groups: []string{},
	}
	return h
}

func (h *TrpcLogHandler) GetLogger(ctx context.Context) log.Logger {
	if l, ok := codec.Message(ctx).Logger().(log.Logger); ok {
		return l
	}
	return h.logger
}

// Enabled implements slog.Handler.
func (h *TrpcLogHandler) Enabled(_ context.Context, level slog.Level) bool {
	// always enabled
	log.EnableTrace()
	return true
}

// Handle implements slog.Handler.
func (t *TrpcLogHandler) Handle(ctx context.Context, record slog.Record) error {

	// 构建键值对
	attrs := make([]log.Field, 0, record.NumAttrs())
	record.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs,
			log.Field{Key: a.Key, Value: a.Value},
		)
		return true
	})

	logger := t.GetLogger(ctx)
	logger = logger.With(attrs...)

	switch record.Level {
	case slog.LevelDebug:
		logger.Debug(record.Message)
	case slog.LevelError:
		logger.Error(record.Message)
	case slog.LevelInfo:
		logger.Info(record.Message)
	case slog.LevelWarn:
		logger.Warn(record.Message)
	default:
		logger.Info(record.Message)
	}
	return nil
}

// WithAttrs implements slog.Handler.
func (t *TrpcLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {

	var t2 = *t

	var newlogger = t.logger
	for _, attr := range attrs {
		newlogger = newlogger.With(t.ToTrpcField(attr))
	}
	t2.logger = newlogger
	return &t2
}

func (t *TrpcLogHandler) ToTrpcField(attr slog.Attr) log.Field {

	return log.Field{
		Key:   attr.Key,
		Value: attr.Value.Any(),
	}
}

// WithGroup implements slog.Handler.
func (t *TrpcLogHandler) WithGroup(name string) slog.Handler {
	var t2 = TrpcLogHandler{}
	t2 = *t
	t2.groups = append(t2.groups, name)
	return &t2
}
