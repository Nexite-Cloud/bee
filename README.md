# BEE

BEE is a lightweight and flexible worker pool implementation in Go, designed to handle concurrent tasks efficiently. It provides configurable options for worker numbers, queue size, and logging, making it suitable for various use cases.

## Features

- Configurable worker pool with default settings.
- Support for custom logging.
- Graceful handling of tasks with error logging.
- Easy-to-use API for pushing tasks and managing workers.

## Installation

To use BEE in your project, add it to your `go.mod` file:

```bash
go get github.com/Nexite-Cloud/bee
```

## Usage

### Basic Example

```go
package main

import (
	"context"
	"fmt"
	"time"
	"github.com/Nexite-Cloud/bee"
)

func main() {
	ctx := context.Background()

	// Create a new Hive with default configuration
	hive := bee.NewHive[int](nil)

	// Set a handler for processing tasks
	hive.SetHandler(func(ctx context.Context, data int) error {
		fmt.Printf("Processing data: %d\n", data)
		return nil
	})

	// Start the worker pool
	hive.Start(ctx)

	// Push tasks into the hive
	for i := 0; i < 10; i++ {
		hive.Push(i)
	}

	// Wait for all tasks to complete
	hive.Wait()
}
```

### Custom Configuration

```go
config := &bee.HiveConfig{
	WorkerNumber: 5,
	QueueSize:    100,
	Logger:       customLogger{},
}

hive := bee.NewHive[int](config)
```

## Configuration Options

- **WorkerNumber**: Number of workers in the pool (default: 1).
- **QueueSize**: Size of the task queue (default: 256).
- **Logger**: Custom logger implementation.

## License

This project is licensed under the MIT License. See the `LICENSE.md` file for details.