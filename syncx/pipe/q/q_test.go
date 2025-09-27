package q

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// =============================================================================
// Common Test Functions using Queue Interface
// =============================================================================

// collectConcurrentPopResults collects results from concurrent pop operations
func collectConcurrentPopResults[T comparable](t *testing.T, q Queue[T]) []T {
	t.Helper()

	results := make(chan T, 2)

	drainFn := func() {
		if val, err := q.Pop(); err == nil {
			results <- val
		}
	}

	go drainFn()
	go drainFn()

	// Collect results with timeout
	var poppedValues []T
	for i := 0; i < 2; i++ {
		select {
		case val := <-results:
			poppedValues = append(poppedValues, val)
		case <-time.After(time.Second):
			t.Error("Pop operation timed out")
		}
	}

	if len(poppedValues) != 2 {
		t.Errorf("Expected 2 values, got %d", len(poppedValues))
	}

	return poppedValues
}

// validatePoppedValues validates that popped values are from the test set
func validatePoppedValues[T comparable](t *testing.T, poppedValues, testItems []T) map[T]bool {
	t.Helper()

	poppedSet := make(map[T]bool)
	for _, poppedVal := range poppedValues {
		found := false
		for _, expectedVal := range testItems {
			if poppedVal == expectedVal {
				found = true
				poppedSet[poppedVal] = true
				break
			}
		}
		if !found {
			t.Errorf("Unexpected value in concurrent pop: %v", poppedVal)
		}
	}
	return poppedSet
}

// popAndValidateRemaining pops remaining items and validates they are unused items
func popAndValidateRemaining[T comparable](t *testing.T, q Queue[T], testItems []T, poppedSet map[T]bool) {
	t.Helper()

	remainingCount := len(testItems) - 2
	for i := 0; i < remainingCount; i++ {
		val, err := q.Pop()
		if err != nil {
			t.Errorf("Failed to pop remaining item %d: %v", i, err)
		}

		found := false
		for _, expectedVal := range testItems {
			if val == expectedVal && !poppedSet[val] {
				found = true
				poppedSet[val] = true
				break
			}
		}
		if !found {
			t.Errorf("Unexpected remaining value: %v", val)
		}
	}
}

// testConcurrentPopValidation tests concurrent pop operations and validates results
func testConcurrentPopValidation[T comparable](t *testing.T, q Queue[T], testItems []T) {
	t.Helper()

	// Collect results from concurrent pops
	poppedValues := collectConcurrentPopResults(t, q)

	// Validate the popped values are from our test set
	poppedSet := validatePoppedValues(t, poppedValues, testItems)

	// Pop and validate remaining items
	popAndValidateRemaining(t, q, testItems, poppedSet)
}

// testBasicOperations tests basic queue operations using the interface
func testBasicOperations[T comparable](t *testing.T, q Queue[T], testItems []T) {
	t.Helper()

	// Test initial state
	if !q.IsEmpty() {
		t.Error("New queue should be empty")
	}
	if q.Len() != 0 {
		t.Errorf("New queue length should be 0, got %d", q.Len())
	}

	// Test capacity behavior for finite capacity queues
	if q.Cap() > 0 && q.IsFull() {
		t.Error("New queue should not be full")
	}

	// Test push
	for i, item := range testItems {
		if err := q.Push(item); err != nil {
			t.Errorf("Failed to push item %d: %v", i, err)
		}
	}

	// Test length after push
	if q.Len() != len(testItems) {
		t.Errorf("Expected length %d, got %d", len(testItems), q.Len())
	}

	// Test concurrent pop pattern (like original tests) - pop 2 items concurrently
	if len(testItems) >= 2 {
		testConcurrentPopValidation(t, q, testItems)
	} else {
		// Fallback: pop in order for small test sets
		for i, expected := range testItems {
			val, err := q.Pop()
			if err != nil {
				t.Errorf("Failed to pop item %d: %v", i, err)
			}
			if val != expected {
				t.Errorf("Expected %v, got %v", expected, val)
			}
		}
	}

	// Queue should be empty now
	if !q.IsEmpty() {
		t.Error("Queue should be empty after popping all items")
	}
}

// testCapacityBehavior tests capacity-related behavior
func testCapacityBehavior[T comparable](t *testing.T, q Queue[T], capacity int, testItem T) {
	t.Helper()

	if capacity <= 0 {
		// Unlimited capacity - should never be full
		if q.IsFull() {
			t.Error("Unlimited capacity queue should never be full")
		}
		return
	}

	// Fixed capacity queue
	if q.Cap() != capacity {
		t.Errorf("Queue capacity should be %d, got %d", capacity, q.Cap())
	}

	// Fill to capacity
	for i := 0; i < capacity; i++ {
		if err := q.Push(testItem); err != nil {
			t.Errorf("Failed to push item %d to capacity: %v", i, err)
		}
	}

	// Should be full now
	if !q.IsFull() {
		t.Error("Queue should be full after pushing to capacity")
	}

	// Next push should fail
	if err := q.Push(testItem); err != ErrQueueFull {
		t.Errorf("Expected ErrQueueFull, got %v", err)
	}
}

// testCloseQueue tests queue closing behavior
func testCloseQueue[T any](t *testing.T, q Queue[T]) {
	t.Helper()

	// Close the queue
	q.Close()

	// Push should fail
	var zero T
	if err := q.Push(zero); err != ErrClosed {
		t.Errorf("Expected ErrClosed when pushing to closed queue, got %v", err)
	}

	// Pop should fail immediately
	done := make(chan bool, 1)
	go func() {
		if _, err := q.Pop(); err != ErrClosed {
			t.Errorf("Expected ErrClosed when popping from closed queue, got %v", err)
		}
		done <- true
	}()

	select {
	case <-done:
		// Good - pop returned immediately
	case <-time.After(100 * time.Millisecond):
		t.Error("Pop should return immediately when queue is closed")
	}

	// Should be marked as closed
	if !q.IsClosed() {
		t.Error("Queue should be marked as closed")
	}
}

// testConcurrentPushPop tests concurrent push/pop operations
func testConcurrentPushPop[T comparable](t *testing.T, q Queue[T], testItems []T) {
	t.Helper()

	var wg sync.WaitGroup
	results := make(chan T, len(testItems))

	// Producer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, item := range testItems {
			for {
				if err := q.Push(item); err == nil {
					break
				}
				time.Sleep(time.Microsecond)
			}
		}
	}()

	// Consumer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < len(testItems); i++ {
			val, err := q.Pop()
			if err != nil {
				t.Errorf("Pop failed: %v", err)
				return
			}
			results <- val
		}
	}()

	wg.Wait()
	close(results)

	// Collect and verify results
	poppedValues := make([]T, 0, len(testItems))
	for val := range results {
		poppedValues = append(poppedValues, val)
	}

	if len(poppedValues) != len(testItems) {
		t.Errorf("Expected %d items, got %d", len(testItems), len(poppedValues))
	}
}

// testReset tests the Reset functionality (if queue supports it)
func testReset[T comparable](t *testing.T, q Queue[T], testItems []T) {
	t.Helper()

	// Add some items
	for _, item := range testItems {
		if err := q.Push(item); err != nil {
			t.Fatalf("Failed to push item in testReset: %v", err)
		}
	}

	// Reset the queue
	q.Reset()

	if !q.IsEmpty() {
		t.Error("Queue should be empty after reset")
	}
	if q.Len() != 0 {
		t.Errorf("Queue length should be 0 after reset, got %d", q.Len())
	}

	// Should be able to push again
	if err := q.Push(testItems[0]); err != nil {
		t.Errorf("Should be able to push after reset: %v", err)
	}
}

// testDetailedConcurrency tests detailed concurrent behavior like original tests
func testDetailedConcurrency[T comparable](t *testing.T, q Queue[T], testItems []T) {
	t.Helper()

	var wg sync.WaitGroup

	// Producer goroutine - similar to original test pattern
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < len(testItems); i++ {
			for {
				if err := q.Push(testItems[i]); err == nil {
					break
				}
				time.Sleep(time.Microsecond)
			}
		}
	}()

	// Consumer goroutine with timeout checking - similar to original test
	wg.Add(1)
	go func() {
		defer wg.Done()
		count := 0
		for count < len(testItems) {
			val, err := q.Pop()
			if err == nil {
				count++
				// Verify the value is one of our test items
				found := false
				for _, expected := range testItems {
					if val == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Unexpected value popped: %v", val)
				}
			} else if err == ErrClosed {
				break
			} else {
				time.Sleep(time.Microsecond)
			}
		}
	}()

	wg.Wait()
}

// testConcurrentPopWithTimeout tests Pop operations with timeout (like original tests)
func testConcurrentPopWithTimeout[T comparable](t *testing.T, q Queue[T], testItems []T) {
	t.Helper()

	// Fill queue first
	for _, item := range testItems {
		if err := q.Push(item); err != nil {
			t.Fatalf("Failed to push test item: %v", err)
		}
	}

	// Test concurrent pop with channel and timeout pattern (from original tests)
	results := make(chan T, len(testItems))
	numPops := 2
	if len(testItems) < 2 {
		numPops = len(testItems)
	}

	for i := 0; i < numPops; i++ {
		go func() {
			if val, err := q.Pop(); err == nil {
				results <- val
			}
		}()
	}

	// Collect results with timeout
	var poppedValues []T
	for i := 0; i < numPops; i++ {
		select {
		case val := <-results:
			poppedValues = append(poppedValues, val)
		case <-time.After(time.Second):
			t.Error("Pop operation timed out")
		}
	}

	// Verify we got expected number of values
	if len(poppedValues) != numPops {
		t.Errorf("Expected %d values, got %d", numPops, len(poppedValues))
	}

	// Verify all popped values are from our test set
	for _, val := range poppedValues {
		found := false
		for _, expected := range testItems {
			if val == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Unexpected value in results: %v", val)
		}
	}
}

// testOriginalInitialState tests the initial state of queue
func testOriginalInitialState(t *testing.T, q Queue[int]) {
	t.Helper()

	if !q.IsEmpty() {
		t.Error("New queue should be empty")
	}
	if q.Cap() > 0 && q.IsFull() {
		t.Error("New queue should not be full")
	}
	if q.Len() != 0 {
		t.Errorf("New queue length should be 0, got %d", q.Len())
	}
}

// testOriginalPushPattern tests pushing items 1, 2, 3
func testOriginalPushPattern(t *testing.T, q Queue[int]) {
	t.Helper()

	if err := q.Push(1); err != nil {
		t.Errorf("Failed to push to empty queue: %v", err)
	}
	if err := q.Push(2); err != nil {
		t.Errorf("Failed to push second item: %v", err)
	}
	if err := q.Push(3); err != nil {
		t.Errorf("Failed to push third item: %v", err)
	}
}

// testOriginalCapacityBehavior tests capacity-limited queue behavior
func testOriginalCapacityBehavior(t *testing.T, q Queue[int]) {
	t.Helper()

	if q.Cap() > 0 && q.Cap() <= 3 {
		if !q.IsFull() {
			t.Error("Queue should be full after pushing to capacity")
		}
		if err := q.Push(4); err != ErrQueueFull {
			t.Errorf("Expected ErrQueueFull, got %v", err)
		}
	}
}

// testOriginalConcurrentPop tests concurrent pop pattern and validates results
func testOriginalConcurrentPop(t *testing.T, q Queue[int]) {
	t.Helper()

	results := make(chan int, 2)
	drainFn := func() {
		if val, err := q.Pop(); err == nil {
			results <- val
		}
	}
	go drainFn()
	go drainFn()

	// Collect results with timeout
	var poppedValues []int
	for i := 0; i < 2; i++ {
		select {
		case val := <-results:
			poppedValues = append(poppedValues, val)
		case <-time.After(time.Second):
			t.Error("Pop operation timed out")
		}
	}

	// Check exact pattern from original tests
	if len(poppedValues) != 2 {
		t.Errorf("Expected 2 values, got %d", len(poppedValues))
	}

	// Verify both values are from {1, 2, 3} and are different
	validValues := map[int]bool{1: true, 2: true, 3: true}
	seenValues := make(map[int]bool)

	for _, v := range poppedValues {
		if !validValues[v] {
			t.Errorf("Unexpected value popped: %d, expected one of {1, 2, 3}", v)
		}
		if seenValues[v] {
			t.Errorf("Duplicate value popped: %d", v)
		}
		seenValues[v] = true
	}
}

// testOriginalFinalOperations tests final push and size check
func testOriginalFinalOperations(t *testing.T, q Queue[int]) {
	t.Helper()

	// Now we can push again (from original test)
	if err := q.Push(4); err != nil {
		t.Errorf("Should be able to push after pop: %v", err)
	}

	// Test size (from original test)
	if q.Len() != 2 {
		t.Errorf("Expected length 2, got %d", q.Len())
	}
}

// testOriginalBasicOperationsPattern replicates the exact pattern from original tests
// This function is specifically for int queues to match original test patterns
func testOriginalBasicOperationsPattern(t *testing.T, q Queue[int]) {
	t.Helper()

	// This function tests the exact pattern from original TestRingQ and TestOriginalQ BasicOperations
	// It assumes the queue has capacity >= 3 and tests with values 1, 2, 3

	testOriginalInitialState(t, q)
	testOriginalPushPattern(t, q)
	testOriginalCapacityBehavior(t, q)
	testOriginalConcurrentPop(t, q)
	testOriginalFinalOperations(t, q)
}

// =============================================================================
// Interface-based Tests - These replace the duplicated tests
// =============================================================================

// TestQueueImplementations tests both queue implementations using interface
func TestQueueImplementations(t *testing.T) {
	testCases := []struct {
		name        string
		createQueue func(capacity int) Queue[int]
	}{
		{
			name: "RingQ",
			createQueue: func(capacity int) Queue[int] {
				if capacity <= 0 {
					capacity = 100 // RingQ requires positive capacity
				}
				return NewRingQ[int](capacity)
			},
		},
		{
			name: "Q",
			createQueue: func(capacity int) Queue[int] {
				return NewQ[int](capacity)
			},
		},
	}

	testItems := []int{1, 2, 3, 4, 5}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("BasicOperations", func(t *testing.T) {
				q := tc.createQueue(10)
				testBasicOperations(t, q, testItems)
			})

			t.Run("OriginalBasicPattern", func(t *testing.T) {
				q := tc.createQueue(5) // Enough capacity for original test pattern
				testOriginalBasicOperationsPattern(t, q)
			})

			t.Run("CapacityBehavior", func(t *testing.T) {
				q := tc.createQueue(3)
				testCapacityBehavior(t, q, 3, 42)
			})

			t.Run("UnlimitedCapacity", func(t *testing.T) {
				q := tc.createQueue(0) // Only Q supports unlimited capacity
				if tc.name == "Q" {
					testCapacityBehavior(t, q, 0, 42)
				}
			})

			t.Run("CloseQueue", func(t *testing.T) {
				q := tc.createQueue(10)
				testCloseQueue(t, q)
			})

			t.Run("ConcurrentOperations", func(t *testing.T) {
				q := tc.createQueue(100)
				testConcurrentPushPop(t, q, testItems)
			})

			t.Run("DetailedConcurrency", func(t *testing.T) {
				q := tc.createQueue(100)
				testDetailedConcurrency(t, q, testItems)
			})

			t.Run("ConcurrentPopWithTimeout", func(t *testing.T) {
				q := tc.createQueue(10)
				testConcurrentPopWithTimeout(t, q, testItems)
			})

			t.Run("Reset", func(t *testing.T) {
				q := tc.createQueue(10)
				testReset(t, q, testItems)
			})
		})
	}
}

// TestQueueGenericTypes tests different types using interface
func TestQueueGenericTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	testCases := []struct {
		name        string
		createQueue func() Queue[Person]
		testItems   []Person
	}{
		{
			name: "RingQ_Person",
			createQueue: func() Queue[Person] {
				return NewRingQ[Person](10)
			},
			testItems: []Person{
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 25},
			},
		},
		{
			name: "Q_Person",
			createQueue: func() Queue[Person] {
				return NewQ[Person](10)
			},
			testItems: []Person{
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 25},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := tc.createQueue()
			testBasicOperations(t, q, tc.testItems)
		})
	}

	// Test string queues
	stringTestCases := []struct {
		name        string
		createQueue func() Queue[string]
	}{
		{
			name:        "RingQ_String",
			createQueue: func() Queue[string] { return NewRingQ[string](10) },
		},
		{
			name:        "Q_String",
			createQueue: func() Queue[string] { return NewQ[string](10) },
		},
	}

	stringItems := []string{"hello", "world", "queue", "test"}
	for _, tc := range stringTestCases {
		t.Run(tc.name, func(t *testing.T) {
			q := tc.createQueue()
			testBasicOperations(t, q, stringItems)
		})
	}

	// Test pointer types
	t.Run("PointerTypes", func(t *testing.T) {
		// Test RingQ with pointers
		t.Run("RingQ_Pointer", func(t *testing.T) {
			q := NewRingQ[*int](5)

			num1, num2, num3 := 42, 84, 126
			items := []*int{&num1, &num2, &num3}

			for _, item := range items {
				if err := q.Push(item); err != nil {
					t.Errorf("Failed to push pointer: %v", err)
				}
			}

			for i, expected := range items {
				val, err := q.Pop()
				if err != nil {
					t.Errorf("Failed to pop pointer %d: %v", i, err)
				}
				if val != expected || *val != *expected {
					t.Errorf("Expected pointer to %d, got pointer to %d", *expected, *val)
				}
			}
		})

		// Test Q with pointers
		t.Run("Q_Pointer", func(t *testing.T) {
			q := NewQ[*int](5)

			num1, num2, num3 := 42, 84, 126
			items := []*int{&num1, &num2, &num3}

			for _, item := range items {
				if err := q.Push(item); err != nil {
					t.Errorf("Failed to push pointer: %v", err)
				}
			}

			for i, expected := range items {
				val, err := q.Pop()
				if err != nil {
					t.Errorf("Failed to pop pointer %d: %v", i, err)
				}
				if val != expected || *val != *expected {
					t.Errorf("Expected pointer to %d, got pointer to %d", *expected, *val)
				}
			}
		})
	})
}

// =============================================================================
// RingQ-specific Tests (tests unique behavior like circular wrapping)
// =============================================================================

func TestRingQSpecific(t *testing.T) {
	t.Run("CircularBehavior", func(t *testing.T) {
		q := NewRingQ[string](2)

		// Fill the queue
		if err := q.Push("a"); err != nil {
			t.Fatalf("Failed to push 'a': %v", err)
		}
		if err := q.Push("b"); err != nil {
			t.Fatalf("Failed to push 'b': %v", err)
		}

		// Pop one item
		val, err := q.Pop()
		if err != nil {
			t.Fatalf("Failed to pop first item: %v", err)
		}
		if val != "a" {
			t.Errorf("Expected 'a', got %v", val)
		}

		// Push another item (should wrap around)
		if err := q.Push("c"); err != nil {
			t.Fatalf("Failed to push 'c': %v", err)
		}

		// Pop remaining items
		val, err = q.Pop()
		if err != nil {
			t.Fatalf("Failed to pop second item: %v", err)
		}
		if val != "b" {
			t.Errorf("Expected 'b', got %v", val)
		}

		val, err = q.Pop()
		if err != nil {
			t.Fatalf("Failed to pop third item: %v", err)
		}
		if val != "c" {
			t.Errorf("Expected 'c', got %v", val)
		}
	})
}

// =============================================================================
// Q-specific Tests (tests unique behavior like unlimited capacity)
// =============================================================================

func TestQSpecific(t *testing.T) {
	t.Run("UnlimitedCapacity", func(t *testing.T) {
		q := NewQ[int](0) // 0 means unlimited

		// Should never be full
		for i := 0; i < 1000; i++ {
			if err := q.Push(i); err != nil {
				t.Errorf("Push failed on unlimited queue: %v", err)
			}
			if q.IsFull() {
				t.Error("Unlimited capacity queue should never be full")
			}
		}

		if q.Len() != 1000 {
			t.Errorf("Expected length 1000, got %d", q.Len())
		}
	})

	t.Run("LimitedCapacity", func(t *testing.T) {
		q := NewQ[string](2)

		// Fill to capacity
		if err := q.Push("a"); err != nil {
			t.Fatalf("Failed to push 'a': %v", err)
		}
		if err := q.Push("b"); err != nil {
			t.Fatalf("Failed to push 'b': %v", err)
		}

		// Should be full
		if !q.IsFull() {
			t.Error("Queue should be full")
		}

		// Next push should fail
		if err := q.Push("c"); err != ErrQueueFull {
			t.Errorf("Expected ErrQueueFull, got %v", err)
		}

		// Pop and push should work
		val, err := q.Pop()
		if err != nil {
			t.Fatalf("Failed to pop: %v", err)
		}
		if val != "a" {
			t.Errorf("Expected 'a', got %v", val)
		}

		if err := q.Push("c"); err != nil {
			t.Errorf("Should be able to push after pop: %v", err)
		}
	})
}

// =============================================================================
// Benchmark Tests using Interface
// =============================================================================

func BenchmarkQueueImplementations(b *testing.B) {
	implementations := []struct {
		name        string
		createQueue func(capacity int) Queue[int]
	}{
		{
			name: "RingQ",
			createQueue: func(capacity int) Queue[int] {
				return NewRingQ[int](capacity)
			},
		},
		{
			name: "Q",
			createQueue: func(capacity int) Queue[int] {
				return NewQ[int](capacity)
			},
		},
	}

	capacities := []int{100, 1000}

	for _, impl := range implementations {
		for _, cap := range capacities {
			b.Run(fmt.Sprintf("%s_Cap%d", impl.name, cap), func(b *testing.B) {
				q := impl.createQueue(cap)
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					if err := q.Push(i % cap); err != nil {
						b.Fatalf("Push failed: %v", err)
					}
					if _, err := q.Pop(); err != nil {
						b.Fatalf("Pop failed: %v", err)
					}
				}
			})
		}
	}
}

func BenchmarkPushOnly(b *testing.B) {
	implementations := []struct {
		name        string
		createQueue func(capacity int) Queue[int]
	}{
		{
			name:        "RingQ",
			createQueue: func(capacity int) Queue[int] { return NewRingQ[int](capacity) },
		},
		{
			name:        "Q",
			createQueue: func(capacity int) Queue[int] { return NewQ[int](capacity) },
		},
	}

	for _, impl := range implementations {
		b.Run(impl.name, func(b *testing.B) {
			q := impl.createQueue(b.N)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				if err := q.Push(i); err != nil {
					b.Fatalf("Push failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkPopOnly(b *testing.B) {
	implementations := []struct {
		name        string
		createQueue func(capacity int) Queue[int]
	}{
		{
			name:        "RingQ",
			createQueue: func(capacity int) Queue[int] { return NewRingQ[int](capacity) },
		},
		{
			name:        "Q",
			createQueue: func(capacity int) Queue[int] { return NewQ[int](capacity) },
		},
	}

	for _, impl := range implementations {
		b.Run(impl.name, func(b *testing.B) {
			q := impl.createQueue(b.N)
			// Pre-fill the queue
			for i := 0; i < b.N; i++ {
				if err := q.Push(i); err != nil {
					b.Fatalf("Push failed during pre-fill: %v", err)
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := q.Pop(); err != nil {
					b.Fatalf("Pop failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkConcurrent(b *testing.B) {
	implementations := []struct {
		name        string
		createQueue func() Queue[int]
	}{
		{
			name:        "RingQ",
			createQueue: func() Queue[int] { return NewRingQ[int](1000) },
		},
		{
			name:        "Q",
			createQueue: func() Queue[int] { return NewQ[int](1000) },
		},
	}

	for _, impl := range implementations {
		b.Run(impl.name, func(b *testing.B) {
			q := impl.createQueue()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					if i%2 == 0 {
						//nolint:errcheck
						q.Push(i)
					} else {
						//nolint:errcheck
						q.Pop()
					}
					i++
				}
			})
		})
	}
}

// BenchmarkDetailedComparison provides detailed performance comparison like original tests
func BenchmarkDetailedComparison(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("RingQ_Cap%d", size), func(b *testing.B) {
			q := NewRingQ[int](size)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				//nolint:errcheck
				q.Push(i % size)
				//nolint:errcheck
				q.Pop()
			}
		})

		b.Run(fmt.Sprintf("Q_Cap%d", size), func(b *testing.B) {
			q := NewQ[int](size)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				//nolint:errcheck
				q.Push(i % size)
				//nolint:errcheck
				q.Pop()
			}
		})
	}
}

// BenchmarkPushOnlyComparison tests pure push performance like original tests
func BenchmarkPushOnlyComparison(b *testing.B) {
	b.Run("RingQ", func(b *testing.B) {
		q := NewRingQ[int](b.N)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			//nolint:errcheck
			q.Push(i)
		}
	})

	b.Run("Q", func(b *testing.B) {
		q := NewQ[int](b.N)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			//nolint:errcheck
			q.Push(i)
		}
	})
}

// BenchmarkPopOnlyComparison tests pure pop performance like original tests
func BenchmarkPopOnlyComparison(b *testing.B) {
	b.Run("RingQ", func(b *testing.B) {
		q := NewRingQ[int](b.N)
		// Pre-fill the queue
		for i := 0; i < b.N; i++ {
			//nolint:errcheck
			q.Push(i)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			//nolint:errcheck
			q.Pop()
		}
	})

	b.Run("Q", func(b *testing.B) {
		q := NewQ[int](b.N)
		// Pre-fill the queue
		for i := 0; i < b.N; i++ {
			//nolint:errcheck
			q.Push(i)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			//nolint:errcheck
			q.Pop()
		}
	})
}

// BenchmarkConcurrentComparison tests concurrent performance like original tests
func BenchmarkConcurrentComparison(b *testing.B) {
	queueSize := 1000

	b.Run("RingQ", func(b *testing.B) {
		q := NewRingQ[int](queueSize)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				if i%2 == 0 {
					//nolint:errcheck
					q.Push(i)
				} else {
					//nolint:errcheck
					q.Pop()
				}
				i++
			}
		})
	})

	b.Run("Q", func(b *testing.B) {
		q := NewQ[int](queueSize)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				if i%2 == 0 {
					//nolint:errcheck
					q.Push(i)
				} else {
					//nolint:errcheck
					q.Pop()
				}
				i++
			}
		})
	})
}
