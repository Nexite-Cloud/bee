package bee

import "context"

type ContextKey string

const (
	CtxKeyWorkerIdx ContextKey = "_worker_idx"
)

// GetWorkerIndex retrieves the worker index from the context. If the index is not set, it returns -1
func GetWorkerIndex(ctx context.Context) int {
	if idx, ok := ctx.Value(CtxKeyWorkerIdx).(int); ok {
		return idx
	}
	return -1 // or any other default value you prefer
}
