package q

import (
	"sync"
	"testing"
	"time"
)

// TestQ_PushBlockingPeek Tests the new PushBlocking and Peek methods for both Q and RingQ implementations
func TestQ_PushBlockingPeek(t *testing.T) {
	implementations := []struct {
		name        string
		createQueue func(capacity int) Queue[int]
	}{
		{
			name: "RingQ",
			createQueue: func(capacity int) Queue[int] {
				if capacity <= 0 {
					capacity = 10 // RingQ requires positive capacity
				}
				return NewRingQ[int](capacity)
			},
		},
		{
			name:        "Q",
			createQueue: func(capacity int) Queue[int] { return NewQ[int](capacity) },
		},
	}

	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			testPeekMethod(t, impl.createQueue)
			testPushBlockingMethod(t, impl.createQueue)
		})
	}
}

func testPeekMethod(t *testing.T, createQueue func(int) Queue[int]) {
	t.Helper()
	t.Run("Peek", func(t *testing.T) {
		q := createQueue(10)

		// Test peek on empty queue
		val := q.Peek()
		if val != 0 {
			t.Errorf("Peek on empty queue should return zero value, got %v", val)
		}

		// Push some items and test peek
		if err := q.Push(42); err != nil {
			t.Fatalf("Failed to push: %v", err)
		}
		if err := q.Push(84); err != nil {
			t.Fatalf("Failed to push: %v", err)
		}

		// Peek should return first item without removing it
		val = q.Peek()
		if val != 42 {
			t.Errorf("Expected peek to return 42, got %v", val)
		}

		// Queue length should be unchanged
		if q.Len() != 2 {
			t.Errorf("Expected queue length 2 after peek, got %d", q.Len())
		}

		// Pop and verify peek returns next item
		poppedVal, err := q.Pop()
		if err != nil {
			t.Fatalf("Failed to pop: %v", err)
		}
		if poppedVal != 42 {
			t.Errorf("Expected pop to return 42, got %v", poppedVal)
		}

		val = q.Peek()
		if val != 84 {
			t.Errorf("Expected peek to return 84 after pop, got %v", val)
		}
	})
}

func testPushBlockingMethod(t *testing.T, createQueue func(int) Queue[int]) {
	t.Helper()
	t.Run("PushBlocking", func(t *testing.T) {
		t.Run("UnlimitedCapacity", func(t *testing.T) {
			if createQueue == nil {
				t.Skip("Skipping unlimited capacity test for RingQ")
				return
			}
			// Only test unlimited capacity for Q
			q := createQueue(0)
			if q.Cap() == 0 { // Only Q supports unlimited capacity
				// Should never block for unlimited capacity
				for i := 0; i < 100; i++ {
					if err := q.PushBlocking(i); err != nil {
						t.Errorf("PushBlocking should never fail with unlimited capacity: %v", err)
					}
				}
				if q.Len() != 100 {
					t.Errorf("Expected 100 items, got %d", q.Len())
				}
			}
		})

		t.Run("LimitedCapacity", func(t *testing.T) {
			q := createQueue(2)

			// Should succeed until capacity is reached
			if err := q.PushBlocking(1); err != nil {
				t.Errorf("PushBlocking should succeed: %v", err)
			}
			if err := q.PushBlocking(2); err != nil {
				t.Errorf("PushBlocking should succeed: %v", err)
			}

			// Now queue should be full
			if !q.IsFull() {
				t.Error("Queue should be full")
			}

			// Test blocking behavior with timeout
			done := make(chan bool, 1)
			go func() {
				err := q.PushBlocking(3) // This should block
				if err != nil {
					t.Errorf("PushBlocking failed: %v", err)
				}
				done <- true
			}()

			// Wait a bit to ensure goroutine is blocked
			select {
			case <-done:
				t.Error("PushBlocking should have blocked on full queue")
			case <-time.After(100 * time.Millisecond):
				// Good, it's blocking
			}

			// Pop one item to make space
			_, err := q.Pop()
			if err != nil {
				t.Fatalf("Failed to pop: %v", err)
			}

			// Now the blocking PushBlocking should succeed
			select {
			case <-done:
				// Good, it succeeded
			case <-time.After(time.Second):
				t.Error("PushBlocking should have succeeded after making space")
			}
		})

		t.Run("ClosedQueue", func(t *testing.T) {
			q := createQueue(10)
			q.Close()

			err := q.PushBlocking(1)
			if err != ErrClosed {
				t.Errorf("Expected ErrClosed, got %v", err)
			}
		})
	})
}

// TestProducerConsumerPerformance tests single producer and single consumer performance
func TestProducerConsumerPerformance(t *testing.T) {
	// Test parameters
	const (
		capacity = 1024
	)

	implementations := []struct {
		name        string
		createQueue func() Queue[int64]
	}{
		{
			name:        "RingQ",
			createQueue: func() Queue[int64] { return NewRingQ[int64](capacity) },
		},
		{
			name:        "Q",
			createQueue: func() Queue[int64] { return NewQ[int64](capacity) },
		},
	}

	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			testSingleProducerConsumer(t, impl.createQueue(), 100000)
		})
	}

	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			testSingleProducerConsumer(t, impl.createQueue(), 1000000)
		})
	}
}

func testSingleProducerConsumer(t *testing.T, q Queue[int64], n int) {
	t.Helper()

	var wg sync.WaitGroup
	// Record start time
	startTime := time.Now()

	// Producer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := int64(0); i < int64(n); i++ {
			err := q.PushBlocking(i)
			if err != nil {
				t.Errorf("Producer failed to push item %d: %v", i, err)
				return
			}
		}
		t.Logf("Producer finished: pushed %d items", n)
	}()

	// Consumer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			val, err := q.Pop()
			if err != nil {
				t.Errorf("Consumer failed to pop item %d: %v", i, err)
				return
			}
			// Verify data correctness (optional)
			expectedVal := int64(i)
			if val != expectedVal {
				t.Errorf("Expected value %d, got %d at position %d", expectedVal, val, i)
				return
			}
		}
		t.Logf("Consumer finished: popped %d items", n)
	}()

	// Wait for both goroutines to complete
	wg.Wait()

	// Calculate total time and average time per operation
	totalTime := time.Since(startTime)
	avgTimePerOp := totalTime / time.Duration(n)

	t.Logf("Performance Results:")
	t.Logf("  Total operations: %d", n)
	t.Logf("  Total time: %v", totalTime)
	t.Logf("  Average time per operation: %v", avgTimePerOp)
	t.Logf("  Operations per second: %.0f", float64(n)/totalTime.Seconds())

	// Verify queue final state
	if !q.IsEmpty() {
		t.Errorf("Expected queue to be empty after test, but it has %d items", q.Len())
	}
}
