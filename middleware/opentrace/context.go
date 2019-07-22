package opentrace

import "context"

type ctxKey int

const (
	skipCreatingKey ctxKey = iota
)

// ContextWithSkipSpanCreating does not create a new span upon request, does not close the span upon response
func ContextWithSkipSpanCreating(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipCreatingKey, true)
}

// SkipSpanCreatingFromContext gets the flag
func SkipSpanCreatingFromContext(ctx context.Context) bool {
	return ctx.Value(skipCreatingKey) != nil
}
