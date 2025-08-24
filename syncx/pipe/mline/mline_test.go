package mline

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pinealctx/neptune/syncx/pipe"
)

func TestMultiLine_BasicFunctionality(t *testing.T) {
	// Create a multi-line with small slot and queue sizes for testing
	mline := NewMultiLine(pipe.WithSlotSize(3), pipe.WithQSize(10))

	// Verify configuration
	if mline.SlotSize() != 3 {
		t.Errorf("Expected slot size 3, got %d", mline.SlotSize())
	}
	if mline.QSize() != 10 {
		t.Errorf("Expected queue size 10, got %d", mline.QSize())
	}

	// Start the multi-line
	mline.Run()
	defer mline.Stop()

	// Test basic async call
	ctx := context.Background()
	callCtx := NewCallCtx(123, func(_ context.Context, sIndex int, req any) (any, error) {
		val, ok := req.(int)
		if !ok {
			return nil, fmt.Errorf("expected int, got %T", req)
		}
		// Include slot index in result to verify proper routing
		return val*2 + sIndex, nil
	}, 42)

	result, err := mline.AsyncCall(ctx, callCtx)
	if err != nil {
		t.Fatalf("AsyncCall failed: %v", err)
	}

	// Verify result includes slot processing
	resultVal, ok := result.(int)
	if !ok {
		t.Fatalf("Expected int result, got %T", result)
	}

	expectedSlotIndex := mline.IndexOf(123)
	expectedResult := 42*2 + expectedSlotIndex
	if resultVal != expectedResult {
		t.Errorf("Expected %d, got %d", expectedResult, resultVal)
	}
}

func TestMultiLine_IndexOf(t *testing.T) {
	mline := NewMultiLine(pipe.WithSlotSize(5))

	testCases := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{1, 1},
		{4, 4},
		{5, 0},  // 5 % 5 = 0
		{7, 2},  // 7 % 5 = 2
		{-3, 3}, // abs(-3) % 5 = 3
		{-7, 2}, // abs(-7) % 5 = 2
	}

	for _, tc := range testCases {
		result := mline.IndexOf(tc.input)
		if result != tc.expected {
			t.Errorf("IndexOf(%d) = %d, expected %d", tc.input, result, tc.expected)
		}
	}
}

func TestMultiLine_ContextCancellation(t *testing.T) {
	mline := NewMultiLine(pipe.WithSlotSize(2), pipe.WithQSize(5))
	mline.Run()
	defer mline.Stop()

	// Test with canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	callCtx := NewCallCtx(456, func(_ context.Context, _ int, _ any) (any, error) {
		time.Sleep(100 * time.Millisecond) // This should not complete
		return "should not reach here", nil
	}, "test")

	result, err := mline.AsyncCall(ctx, callCtx)
	if err == nil {
		t.Error("Expected context cancellation error")
	}
	if result != nil {
		t.Errorf("Expected nil result on cancellation, got %v", result)
	}
}

func TestMultiLine_ConcurrentOperations(t *testing.T) {
	const (
		slotSize = 4
		numOps   = 100
	)

	mline := NewMultiLine(pipe.WithSlotSize(slotSize), pipe.WithQSize(50))
	mline.Run()
	defer mline.Stop()

	// Counter for each slot to verify load distribution
	var slotCounters [slotSize]int64

	// Launch concurrent operations
	var wg sync.WaitGroup
	results := make(chan int, numOps)
	errors := make(chan error, numOps)

	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()

			ctx := context.Background()
			callCtx := NewCallCtx(val, func(_ context.Context, sIndex int, req any) (any, error) {
				input, ok := req.(int)
				if !ok {
					return nil, fmt.Errorf("expected int, got %T", req)
				}

				// Count operations per slot
				atomic.AddInt64(&slotCounters[sIndex], 1)

				// Simulate some work
				time.Sleep(time.Millisecond)

				return input + 1000, nil
			}, val)

			result, err := mline.AsyncCall(ctx, callCtx)
			if err != nil {
				errors <- err
				return
			}
			resultVal, ok := result.(int)
			if !ok {
				errors <- fmt.Errorf("expected int result, got %T", result)
				return
			}
			results <- resultVal
		}(i)
	}

	wg.Wait()
	close(results)
	close(errors)

	// Check for errors
	for err := range errors {
		t.Fatalf("Concurrent operation failed: %v", err)
	}

	// Collect and verify results
	resultCount := 0
	for result := range results {
		if result < 1000 || result >= 1000+numOps {
			t.Errorf("Result %d is not in expected range [1000, %d)", result, 1000+numOps)
		}
		resultCount++
	}

	if resultCount != numOps {
		t.Errorf("Expected %d results, got %d", numOps, resultCount)
	}

	// Verify load distribution (each slot should have some operations)
	for i, count := range slotCounters {
		t.Logf("Slot %d processed %d operations", i, count)
		if count == 0 {
			t.Errorf("Slot %d processed no operations - poor load distribution", i)
		}
	}
}

func TestMultiLine_ErrorHandling(t *testing.T) {
	mline := NewMultiLine(pipe.WithSlotSize(2), pipe.WithQSize(5))
	mline.Run()
	defer mline.Stop()

	ctx := context.Background()
	callCtx := NewCallCtx(789, func(_ context.Context, _ int, _ any) (any, error) {
		return nil, fmt.Errorf("simulated error")
	}, "test")

	result, err := mline.AsyncCall(ctx, callCtx)
	if err == nil {
		t.Error("Expected error from async call")
	}
	if result != nil {
		t.Errorf("Expected nil result on error, got %v", result)
	}

	if err.Error() != "simulated error" {
		t.Errorf("Expected 'simulated error', got '%s'", err.Error())
	}
}

func TestMultiLine_WaitStop(t *testing.T) {
	mline := NewMultiLine(pipe.WithSlotSize(2), pipe.WithQSize(5))
	mline.Run()

	// Test WaitStop with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Stop the multi-line
	go func() {
		time.Sleep(50 * time.Millisecond)
		mline.Stop()
	}()

	// Should complete before timeout
	err := mline.WaitStop(ctx)
	if err != nil {
		t.Errorf("WaitStop failed: %v", err)
	}
}

func TestMultiLine_WaitStopTimeout(t *testing.T) {
	mline := NewMultiLine(pipe.WithSlotSize(2), pipe.WithQSize(5))
	mline.Run()
	defer mline.Stop()

	// Test WaitStop timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Don't stop the multi-line, should timeout
	err := mline.WaitStop(ctx)
	if err == nil {
		t.Error("Expected timeout error")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded, got %v", err)
	}
}

func TestMultiLine_SlotLoadBalancing(t *testing.T) {
	const slotSize = 3
	mline := NewMultiLine(pipe.WithSlotSize(slotSize), pipe.WithQSize(20))
	mline.Run()
	defer mline.Stop()

	// Test that same hash index goes to same slot consistently
	testHash := 12345
	expectedSlot := mline.IndexOf(testHash)

	// Use channel to collect results safely
	resultChan := make(chan int, 10)
	var wg sync.WaitGroup

	// Run multiple operations with same hash
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx := context.Background()
			callCtx := NewCallCtx(testHash, func(_ context.Context, sIndex int, _ any) (any, error) {
				return sIndex, nil
			}, nil)

			result, err := mline.AsyncCall(ctx, callCtx)
			if err != nil {
				t.Errorf("AsyncCall failed: %v", err)
				return
			}

			slotIndex, ok := result.(int)
			if !ok {
				t.Errorf("Expected int result, got %T", result)
				return
			}
			resultChan <- slotIndex
		}()
	}

	wg.Wait()
	close(resultChan)

	// Collect results from channel
	var resultSlots []int
	for slot := range resultChan {
		resultSlots = append(resultSlots, slot)
	}

	// All operations with same hash should go to same slot
	for _, slot := range resultSlots {
		if slot != expectedSlot {
			t.Errorf("Expected all operations to go to slot %d, but got slot %d", expectedSlot, slot)
		}
	}
}

// Benchmark tests
func BenchmarkMultiLine_AsyncCall(b *testing.B) {
	mline := NewMultiLine(pipe.WithSlotSize(4), pipe.WithQSize(1000))
	mline.Run()
	defer mline.Stop()

	callFn := func(_ context.Context, _ int, req any) (any, error) {
		val, ok := req.(int)
		if !ok {
			return nil, fmt.Errorf("expected int, got %T", req)
		}
		return val * 2, nil
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ctx := context.Background()
			callCtx := NewCallCtx(i, callFn, i)
			_, err := mline.AsyncCall(ctx, callCtx)
			if err != nil {
				b.Fatalf("AsyncCall failed: %v", err)
			}
			i++
		}
	})
}

func BenchmarkMultiLine_IndexOf(b *testing.B) {
	mline := NewMultiLine(pipe.WithSlotSize(509)) // Prime number

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mline.IndexOf(i)
	}
}
