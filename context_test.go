package bee

import (
	"context"
	"testing"
)

func TestGetWorkerIndex(t *testing.T) {
	ctx := context.Background()
	idx := GetWorkerIndex(ctx)
	if idx != -1 {
		t.Errorf("expected -1, got %d", idx)
	}
	ctx = context.WithValue(ctx, CtxKeyWorkerIdx, 5)
	idx = GetWorkerIndex(ctx)
	if idx != 5 {
		t.Errorf("expected 5, got %d", idx)
	}
}
