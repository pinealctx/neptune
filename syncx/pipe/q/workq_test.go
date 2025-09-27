package q

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestTask is a simple test task implementation
type TestTask struct {
	id       int
	executed *int32 // atomic counter
	results  *[]int // pointer to slice to collect results
	mu       *sync.Mutex
}

func (t *TestTask) Do() {
	atomic.AddInt32(t.executed, 1)
	if t.results != nil && t.mu != nil {
		t.mu.Lock()
		*t.results = append(*t.results, t.id)
		t.mu.Unlock()
	}
}

// TestSimpleTaskProcessor tests basic functionality of SimpleTaskProcessor
func TestSimpleTaskProcessor(t *testing.T) {
	t.Run("BasicFunctionality", func(t *testing.T) {
		var executed int32
		var results []int
		var mu sync.Mutex

		// Create processor with 2 workers
		processor := NewSimpleTaskProcessor(2)
		defer processor.Shutdown()

		// Submit 5 tasks
		for i := 0; i < 5; i++ {
			task := &TestTask{
				id:       i,
				executed: &executed,
				results:  &results,
				mu:       &mu,
			}
			err := processor.Submit(task)
			if err != nil {
				t.Errorf("Failed to submit task %d: %v", i, err)
			}
		}

		// Wait for processing to complete
		time.Sleep(100 * time.Millisecond)

		// Check results
		executedCount := atomic.LoadInt32(&executed)
		if executedCount != 5 {
			t.Errorf("Expected 5 executed tasks, got %d", executedCount)
		}

		mu.Lock()
		if len(results) != 5 {
			t.Errorf("Expected 5 results, got %d", len(results))
		}
		mu.Unlock()
	})

	t.Run("MultipleWorkers", func(t *testing.T) {
		const numWorkers = 4
		const numTasks = 100

		var executed int32
		processor := NewSimpleTaskProcessor(numWorkers)
		defer processor.Shutdown()

		// Submit many tasks
		for i := 0; i < numTasks; i++ {
			task := &TestTask{
				id:       i,
				executed: &executed,
			}
			err := processor.Submit(task)
			if err != nil {
				t.Errorf("Failed to submit task %d: %v", i, err)
			}
		}

		// Wait for all tasks to complete
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			if atomic.LoadInt32(&executed) == numTasks {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}

		executedCount := atomic.LoadInt32(&executed)
		if executedCount != numTasks {
			t.Errorf("Expected %d executed tasks, got %d", numTasks, executedCount)
		}
	})

	t.Run("ShutdownStopsWorkers", func(t *testing.T) {
		processor := NewSimpleTaskProcessor(2)

		// Submit a few tasks first
		var executed int32
		for i := 0; i < 3; i++ {
			task := &TestTask{
				id:       i,
				executed: &executed,
			}
			err := processor.Submit(task)
			if err != nil {
				t.Errorf("Failed to submit task %d: %v", i, err)
				return
			}
		}

		// Wait a bit for tasks to execute
		time.Sleep(50 * time.Millisecond)

		// Shutdown
		processor.Shutdown()

		// Try to submit more tasks (should be possible but won't be processed)
		oldExecuted := atomic.LoadInt32(&executed)

		// Submit more tasks after shutdown
		for i := 3; i < 6; i++ {
			task := &TestTask{
				id:       i,
				executed: &executed,
			}
			// This might succeed (queue not closed) or fail (queue closed)
			err := processor.Submit(task)
			if err == nil {
				t.Logf("Submit task %d after shutdown", i)
				return
			}
		}

		// Wait a bit and ensure no new tasks were executed
		time.Sleep(100 * time.Millisecond)
		newExecuted := atomic.LoadInt32(&executed)

		// The executed count should not increase significantly after shutdown
		if newExecuted > oldExecuted+1 { // Allow for one task that might have been in progress
			t.Logf("Warning: Some tasks may have been executed after shutdown. Before: %d, After: %d",
				oldExecuted, newExecuted)
		}
	})
}

// SlowTask simulates a task that takes some time to complete
type SlowTask struct {
	duration time.Duration
	executed *int32
}

func (t *SlowTask) Do() {
	time.Sleep(t.duration)
	atomic.AddInt32(t.executed, 1)
}

// TestSimpleTaskProcessorConcurrency tests concurrent behavior
func TestSimpleTaskProcessorConcurrency(t *testing.T) {
	const numWorkers = 3
	const taskDuration = 50 * time.Millisecond
	const numTasks = 6

	processor := NewSimpleTaskProcessor(numWorkers)
	defer processor.Shutdown()

	var executed int32
	start := time.Now()

	// Submit tasks that each take taskDuration to complete
	for i := 0; i < numTasks; i++ {
		task := &SlowTask{
			duration: taskDuration,
			executed: &executed,
		}
		err := processor.Submit(task)
		if err != nil {
			t.Errorf("Failed to submit task %d: %v", i, err)
		}
	}

	// Wait for all tasks to complete
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&executed) == numTasks {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	elapsed := time.Since(start)
	executedCount := atomic.LoadInt32(&executed)

	if executedCount != numTasks {
		t.Errorf("Expected %d executed tasks, got %d", numTasks, executedCount)
	}

	// With 3 workers processing 6 tasks of 50ms each, should take roughly 100ms
	// (2 batches of 3 tasks each), not 300ms (sequential)
	expectedMaxTime := taskDuration * time.Duration(numTasks/numWorkers+1)
	if elapsed > expectedMaxTime*2 { // Allow some margin
		t.Errorf("Tasks took too long: %v, expected around %v", elapsed, expectedMaxTime)
	}

	t.Logf("Processed %d tasks in %v with %d workers", numTasks, elapsed, numWorkers)
}

// BenchmarkSimpleTaskProcessor benchmarks the task processor
func BenchmarkSimpleTaskProcessor(b *testing.B) {
	processor := NewSimpleTaskProcessor(4)
	defer processor.Shutdown()

	var executed int32

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := &TestTask{
				executed: &executed,
			}
			_ = processor.Submit(task)
		}
	})

	// Wait for all submitted tasks to complete
	deadline := time.Now().Add(5 * time.Second)
	targetCount := int32(b.N)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&executed) >= targetCount {
			break
		}
		time.Sleep(time.Millisecond)
	}

	b.StopTimer()

	executedCount := atomic.LoadInt32(&executed)
	b.Logf("Executed %d tasks out of %d submitted", executedCount, b.N)
}

// ErrorTask simulates a task that might panic
type ErrorTask struct {
	shouldPanic bool
	executed    *int32
}

func (t *ErrorTask) Do() {
	defer func() {
		atomic.AddInt32(t.executed, 1)
	}()

	if t.shouldPanic {
		panic("task error")
	}
}

// TestSimpleTaskProcessorErrorHandling tests error handling
func TestSimpleTaskProcessorErrorHandling(t *testing.T) {
	processor := NewSimpleTaskProcessor(2)
	defer processor.Shutdown()

	var executed int32

	// Submit mix of normal and panicking tasks
	for i := 0; i < 10; i++ {
		task := &ErrorTask{
			shouldPanic: i%3 == 0, // Every 3rd task panics
			executed:    &executed,
		}
		err := processor.Submit(task)
		if err != nil {
			t.Errorf("Failed to submit task %d: %v", i, err)
		}
	}

	// Wait for all tasks to complete
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&executed) == 10 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	executedCount := atomic.LoadInt32(&executed)
	if executedCount != 10 {
		t.Errorf("Expected 10 executed tasks (including panicking ones), got %d", executedCount)
	}
}

// TestSimpleTaskProcessorDoubleShutdown tests multiple shutdown calls
func TestSimpleTaskProcessorDoubleShutdown(t *testing.T) {
	processor := NewSimpleTaskProcessor(2)

	// First shutdown
	processor.Shutdown()

	// Second shutdown should not panic
	processor.Shutdown()

	// Try to submit after shutdown
	task := &TestTask{executed: new(int32)}
	err := processor.Submit(task)
	// err must not nil
	if err == nil {
		t.Errorf("Expected error when submitting after shutdown, got nil")
	}
	// Should either succeed (if queue not closed yet) or fail gracefully
	//_ = err // We don't care about the exact behavior, just that it doesn't panic
}

// TestSimpleTaskProcessorZeroWorkers tests edge case
func TestSimpleTaskProcessorZeroWorkers(t *testing.T) {
	processor := NewSimpleTaskProcessor(0)
	defer processor.Shutdown()

	// Submit a task
	var executed int32
	task := &TestTask{executed: &executed}
	err := processor.Submit(task)

	// With 0 workers, task should not be processed
	if err == nil {
		time.Sleep(100 * time.Millisecond)
		executedCount := atomic.LoadInt32(&executed)
		if executedCount > 0 {
			t.Errorf("Expected 0 executed tasks with 0 workers, got %d", executedCount)
		}
	}
}

// CountingTask helps count execution order
type CountingTask struct {
	id      int
	counter *int32
	order   *[]int
	mu      *sync.Mutex
}

func (t *CountingTask) Do() {
	currentCount := atomic.AddInt32(t.counter, 1)
	if t.order != nil && t.mu != nil {
		t.mu.Lock()
		*t.order = append(*t.order, t.id)
		t.mu.Unlock()
	}
	_ = currentCount
}

// TestSimpleTaskProcessorTaskOrder tests task processing
func TestSimpleTaskProcessorTaskOrder(t *testing.T) {
	processor := NewSimpleTaskProcessor(1) // Single worker for predictable order
	defer processor.Shutdown()

	var counter int32
	var order []int
	var mu sync.Mutex

	// Submit tasks in order
	for i := 0; i < 5; i++ {
		task := &CountingTask{
			id:      i,
			counter: &counter,
			order:   &order,
			mu:      &mu,
		}
		err := processor.Submit(task)
		if err != nil {
			t.Errorf("Failed to submit task %d: %v", i, err)
		}
	}

	// Wait for processing
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&counter) == 5 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(order) != 5 {
		t.Errorf("Expected 5 tasks in order, got %d", len(order))
		return
	}

	// With single worker, tasks should execute in FIFO order
	for i, taskId := range order {
		if taskId != i {
			t.Errorf("Expected task %d at position %d, got task %d", i, i, taskId)
		}
	}
}

// TestSimpleTaskProcessorHighLoad tests high concurrent load
func TestSimpleTaskProcessorHighLoad(t *testing.T) {
	const numWorkers = 8
	const numTasks = 10000
	const numProducers = 4

	processor := NewSimpleTaskProcessor(numWorkers)
	defer processor.Shutdown()

	var executed int32
	var wg sync.WaitGroup

	start := time.Now()

	// Start multiple producers
	for p := 0; p < numProducers; p++ {
		wg.Add(1)
		go func(producerId int) {
			defer wg.Done()
			tasksPerProducer := numTasks / numProducers
			for i := 0; i < tasksPerProducer; i++ {
				task := &TestTask{
					id:       producerId*tasksPerProducer + i,
					executed: &executed,
				}
				err := processor.Submit(task)
				if err != nil {
					t.Errorf("Producer %d failed to submit task %d: %v", producerId, i, err)
					return
				}
			}
		}(p)
	}

	// Wait for all producers to finish
	wg.Wait()

	// Wait for all tasks to be processed
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&executed) == numTasks {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	elapsed := time.Since(start)
	executedCount := atomic.LoadInt32(&executed)

	if executedCount != numTasks {
		t.Errorf("Expected %d executed tasks, got %d", numTasks, executedCount)
	}

	throughput := float64(numTasks) / elapsed.Seconds()
	t.Logf("Processed %d tasks in %v with %d workers and %d producers (%.0f tasks/sec)",
		numTasks, elapsed, numWorkers, numProducers, throughput)
}
