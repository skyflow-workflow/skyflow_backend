package trpclog

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"trpc.group/trpc-go/trpc-go/log"
)

func TestSlog(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceid", "traceid-xxxx-xxxx-xxxx")
	slog.Info("Standard log ")
	slog.Info("my first slog msg", "greeting", "hello, slog")
	slog.Error("my secord slog message", "greeting", "hello slog")
	slog.Warn("my third message", "greeting", "hello slog")

	h1 := slog.NewTextHandler(os.Stdout, nil)
	l1 := slog.New(h1)
	l1.Info("Text Logger ")
	l1.Info("greeting", "name", "name l1")
	l1.WithGroup("Group1").WithGroup("Group2").Info("with group info log", "name", "name1")
	l1.DebugContext(ctx, "debug context log with context", "att1", "val1")
	l1.InfoContext(ctx, "info context log with context", "att1", "val1")

	h2 := slog.NewJSONHandler(os.Stdout, nil)
	l2 := slog.New(h2)
	l2.Info("Text Logger ")
	l2.Info("greeting", "name", "xxx")

	l3 := l2.WithGroup("Group1").WithGroup("Group2")
	l3.Info("with group info log", "name", "xxx")
	l4 := l3.With("att1", "val1")
	l4.Info("l4 info with group  attr info log", "name", "xxx")
	l4.Info("l4 info with group  attr info log", "name", "xxx", "att1", "val2", "att3", "val3")

}

func TestTrpcHandler(t *testing.T) {

	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceid", "traceid-xxxx-xxxx-xxxx")

	trpcHandler1 := NewHandlerFromTrpcLogger(nil)
	l1 := slog.New(trpcHandler1)
	l1.Info("Test Trpc Logger")
	l1.Info("my first slog msg", "greeting", "hello, slog")
	l2 := l1.With("att1", "val1").With("att2", "val2")
	l2.Info("l2 with attr info log", "name", "name1")
	l2.Info("l1 with attr info log", "name", "name2")

	var jsonLoggerConfig = []log.OutputConfig{
		{
			Writer:    "console",
			Level:     "debug",
			Formatter: "json",
		},
	}

	trpcLogger2 := log.NewZapLog(jsonLoggerConfig)
	trpcLogger2 = trpcLogger2.With(log.Field{Key: "traceid", Value: "traceid-xxxx-xxxx-xxxx"},
		log.Field{Key: "instance_id", Value: "instance-xxxx-xxxx"},
		log.Field{Key: "service_name", Value: "myservice"},
		log.Field{Key: "service_version", Value: "v1.0.0"},
	)
	ctx = ContextWithLogger(ctx, trpcLogger2)

	l3 := slog.New(NewHandlerFromTrpcLogger(trpcLogger2))
	l3.Info("Test Trpc Logger2")
	l3.Info("my first trpc msg", "greeting", "hello, trpc log")
	l4 := l3.With("att1", "val1").With("att2", "val2")
	l4.Info("l4 with attr info log", "name", "name1")
	l4.Info("l4 with attr info log", "name", "name2")
	l4.DebugContext(ctx, "l4 debug context log with context", "att3", "val3")
}
