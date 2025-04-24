package decoder

import (
	"context"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestContext(t *testing.T) {
	t.Log("TestContext")
	newPath := []string{"path1"}
	newCtx := context.Background()
	ctx := AddPath(newCtx, newPath...)
	ctx1 := Func1(t, ctx)
	assert.Equal(t, GetPath(ctx1), []string{"path1", "func1"})
	ctx2 := Func2(t, ctx)
	assert.Equal(t, GetPath(ctx2), []string{"path1", "func2"})
	ctx2 = Func2(t, ctx1)
	assert.Equal(t, GetPath(ctx2), []string{"path1", "func1", "func2"})
}

func Func1(t *testing.T, ctx context.Context) context.Context {
	ctx = AddPath(ctx, "func1")
	t.Log(GetPath(ctx))
	return ctx
}

func Func2(t *testing.T, ctx context.Context) context.Context {
	ctx = AddPath(ctx, "func2")
	t.Log(GetPath(ctx))
	return ctx
}
