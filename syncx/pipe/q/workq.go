package q

import (
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
)

// ITaskItem represents a task that can be executed by the SimpleTaskProcessor.
// Implementations should ensure that the Do method is safe to call concurrently
// and handles any errors internally, as no error is returned.
type ITaskItem interface {
	// Do executes the task. This method should be idempotent and thread-safe.
	// Any errors that occur during execution should be handled internally.
	Do()
}

// SimpleTaskProcessor is a high-performance, concurrent task processor that uses
// an unbounded queue to manage task execution. It provides a simple interface
// for submitting tasks and automatically distributes them across multiple worker
// goroutines for parallel execution.
//
// Key features:
//   - Unbounded task queue (no capacity limit)
//   - Configurable number of worker goroutines
//   - Automatic task distribution and load balancing
//   - Graceful shutdown with proper cleanup
//   - Thread-safe operations
//
// Example usage:
//
//	type MyTask struct {
//		data string
//	}
//
//	func (t *MyTask) Do() {
//		fmt.Println("Processing:", t.data)
//	}
//
//	processor := NewSimpleTaskProcessor(4) // 4 workers
//	processor.Submit(&MyTask{data: "hello"})
//	processor.Shutdown() // Clean shutdown
//
// The processor is safe for concurrent use by multiple goroutines.
type SimpleTaskProcessor struct {
	// queue is an unbounded queue for storing tasks
	queue Queue[ITaskItem]
	// workerCount is the number of worker goroutines
	workerCount int
}

// NewSimpleTaskProcessor creates and starts a new task processor with the specified
// number of worker goroutines.
//
// Parameters:
//   - workerCount: Number of worker goroutines to spawn. Must be >= 0.
//     If 0, tasks will be queued but not processed until shutdown.
//
// Returns:
//   - *SimpleTaskProcessor: A running task processor ready to accept tasks.
//
// The processor starts immediately and begins listening for tasks. Worker goroutines
// will block waiting for tasks until the processor is shut down.
func NewSimpleTaskProcessor(workerCount int) *SimpleTaskProcessor {
	x := &SimpleTaskProcessor{
		queue:       NewQ[ITaskItem](0),
		workerCount: workerCount,
	}
	x.start()
	return x
}

// Submit adds a task to the processing queue. The task will be executed
// asynchronously by one of the available worker goroutines.
//
// Parameters:
//   - task: The task to be executed. Must implement ITaskItem interface.
//
// Returns:
//   - error: ErrClosed if the processor has been shut down, nil otherwise.
//
// This method is thread-safe and can be called concurrently by multiple goroutines.
// Tasks are processed in FIFO order by the available workers.
func (x *SimpleTaskProcessor) Submit(task ITaskItem) error {
	return x.queue.Push(task)
}

// Shutdown gracefully stops the task processor. It closes the task queue,
// which signals all worker goroutines to finish processing their current tasks
// and then exit.
//
// After calling Shutdown:
//   - No new tasks can be submitted (Submit will return ErrClosed)
//   - Workers will finish processing any tasks already in progress
//   - Workers will exit after processing remaining queued tasks
//   - The processor cannot be restarted
//
// This method is thread-safe and can be called multiple times safely.
// Subsequent calls to Shutdown have no effect.
func (x *SimpleTaskProcessor) Shutdown() {
	x.queue.Close()
}

// start initializes and starts the specified number of worker goroutines.
// Each worker runs in its own goroutine and processes tasks from the queue
// until the queue is closed.
func (x *SimpleTaskProcessor) start() {
	for i := 0; i < x.workerCount; i++ {
		go x.runWorker()
	}
}

// runWorker is the main loop for a worker goroutine. It continuously pops tasks
// from the queue and executes them until the queue is closed (returns ErrClosed).
// This method handles the worker lifecycle and ensures clean shutdown.
// It recovers from panics to prevent a single failing task from crashing the worker.
func (x *SimpleTaskProcessor) runWorker() {
	for {
		task, err := x.queue.Pop()
		if err != nil {
			return
		}

		// Execute task with panic recovery
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Log panic or handle it as needed
					ulog.Error("SimpleTaskProcessor worker panic", zap.Any("panic", r), zap.Stack("stack"))
				}
			}()
			task.Do()
		}()
	}
}
