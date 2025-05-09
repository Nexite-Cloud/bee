package bee

import (
	"context"
	"sync"
)

const DefaultWorkerNumber = 1
const DefaultQueueSize = 1 << 8

// HiveConfig is the config for the hive
type HiveConfig struct {
	WorkerNumber int
	QueueSize    int
	Logger       Logger
}

// NewConfig create a new config with default values
func NewConfig() *HiveConfig {
	return &HiveConfig{
		WorkerNumber: DefaultWorkerNumber,
		QueueSize:    DefaultQueueSize,
		Logger:       NoLog{},
	}
}

// WithWorkerNumber set the number of workers, default is 1
func (h *HiveConfig) WithWorkerNumber(workerNumber int) *HiveConfig {
	h.WorkerNumber = workerNumber
	return h
}

// WithQueueSize set the size of the queue, default is 256
func (h *HiveConfig) WithQueueSize(queueSize int) *HiveConfig {
	h.QueueSize = queueSize
	return h
}

// WithLogger set the logger, default is NoLog
func (h *HiveConfig) WithLogger(logger Logger) *HiveConfig {
	h.Logger = logger
	return h
}

// Handler is the function type for handling data in the hive
type Handler[T any] func(ctx context.Context, data T) error

// Hive is a worker pool that can handle data concurrently
type Hive[T any] struct {
	once    sync.Once
	config  *HiveConfig
	wg      sync.WaitGroup
	cData   chan T
	handler Handler[T]
}

// NewHive create a worker pool with a given config, if config is nil, use NewConfig()
func NewHive[T any](config *HiveConfig) *Hive[T] {
	if config == nil {
		config = NewConfig()
	}
	return &Hive[T]{
		config: config,
		cData:  make(chan T, config.QueueSize),
	}
}

// SetHandler set the handler for the hive, if not set, it will panic when push data into the hive
func (h *Hive[T]) SetHandler(handler Handler[T]) {
	h.handler = handler
}

func (h *Hive[T]) handle(ctx context.Context, workerID int, data T) {
	defer h.wg.Done()
	if err := h.handler(ctx, data); err != nil {
		h.config.Logger.Error(ctx, "handle msg failed", "worker_id", workerID, "data", data, "err", err)
	}
}

// Push data into the hive, if the hive is closed, it will panic
func (h *Hive[T]) Push(data T) {
	h.wg.Add(1)
	h.cData <- data
}

// Start the worker pool, if already started, it will do nothing
func (h *Hive[T]) Start(ctx context.Context) {
	h.once.Do(func() {
		for i := 0; i < h.config.WorkerNumber; i++ {
			go func(id int) {
				h.config.Logger.Info(ctx, "worker started", "worker_id", id)
				for {
					select {
					case <-ctx.Done():
						return
					case data, ok := <-h.cData:
						if !ok {
							return
						}
						h.config.Logger.Info(ctx, "receive data", "worker_id", id, "data", data)
						h.handle(ctx, id, data)
					}
				}
			}(i)
		}
	})
}

// Wait for all workers to finish, if already waited, it will do nothing
func (h *Hive[T]) Wait() {
	h.wg.Wait()
	h.once.Do(func() {
		close(h.cData)
	})
}
