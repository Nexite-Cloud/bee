package bee

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"os"
	"sync"
	"testing"
)

type atomicMap[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

func (am *atomicMap[K, V]) Set(key K, value V) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.m[key] = value
}

func (am *atomicMap[K, V]) Get(key K) (V, bool) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	value, ok := am.m[key]
	return value, ok
}

func (am *atomicMap[K, V]) Data() map[K]V {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.m
}

func TestHive(t *testing.T) {
	logger := NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
	mp := &atomicMap[int, bool]{m: make(map[int]bool)}
	n := int(rand.Int64N(100000))
	h := NewHive[int](NewConfig().WithWorkerNumber(10).WithLogger(logger))
	h.SetHandler(func(_ context.Context, i int) error {
		mp.Set(i, true)
		return nil
	})
	h.Start(t.Context())
	for i := 0; i < n; i++ {
		h.Push(i)
	}
	h.Wait()
	if len(mp.Data()) != n {
		t.Fatal("len(mp.Data()) != n")
	}
}
