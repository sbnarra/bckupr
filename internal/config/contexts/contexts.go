package contexts

import (
	"context"
	"runtime"
)

type key int

const debugKey key = 1
const threadLimitKey key = 2
const nameKey key = 3

func Using(
	ctx context.Context,
	name string,
	debug bool,
	threadLimit int,
) context.Context {
	ctx = WithName(ctx, name)
	ctx = context.WithValue(ctx, debugKey, debug)
	ctx = context.WithValue(ctx, threadLimitKey, threadLimit)
	return ctx
}

func WithName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, nameKey, name)
}

func Name(ctx context.Context) string {
	return value(ctx, nameKey, "")
}

func Debug(ctx context.Context) bool {
	return value(ctx, debugKey, false)
}

func ThreadLimit(ctx context.Context) int {
	return value(ctx, threadLimitKey, runtime.NumCPU())
}

func value[T any](ctx context.Context, key key, fallback T) T {
	if val := ctx.Value(key); val != nil {
		if val, ok := val.(T); ok {
			return val
		}
	}
	return fallback
}
