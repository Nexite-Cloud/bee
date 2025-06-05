package bee

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
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
	n := int(rand.Int63n(100000))
	h := NewHive[int](NewConfig().WithWorkerNumber(10).WithLogger(logger))
	h.SetHandler(func(ctx context.Context, i int) error {
		fmt.Println("worker", GetWorkerIndex(ctx), "processing", i)
		mp.Set(i, true)
		return nil
	})
	h.Start(context.Background())
	for i := 0; i < n; i++ {
		h.Push(i)
	}
	h.Wait()
	if len(mp.Data()) != n {
		t.Fatal("len(mp.Data()) != n")
	}
}

func TestHiveConfig(t *testing.T) {
	logger := NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
	hc := NewConfig().WithWorkerNumber(10).WithQueueSize(100).WithLogger(logger)
	if hc.WorkerNumber != 10 {
		t.Fatal("hc.WorkerNumber != 10")
	}
	if hc.QueueSize != 100 {
		t.Fatal("hc.QueueSize != 100")
	}
	if hc.Logger == nil {
		t.Fatal("hc.Logger == nil")
	}
}
