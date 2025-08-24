package line

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLine_BasicFunctionality(t *testing.T) {
	var wg sync.WaitGroup

	// Create a line with a small queue size
	line := NewLine(&wg, WithQSize(10), WithName("test-line"))

	// Start the line
	line.Run()
	defer line.Stop()

	// Test basic async call
	ctx := context.Background()
	callCtx := NewCallCtx(func(_ context.Context, req any) (any, error) {
		val, ok := req.(int)
		if !ok {
			return nil, fmt.Errorf("expected int, got %T", req)
		}
		return val * 2, nil
	}, 42)

	result, err := line.AsyncCall(ctx, callCtx)
	if err != nil {
		t.Fatalf("AsyncCall failed: %v", err)
	}

	resultInt, ok := result.(int)
	if !ok {
		t.Fatalf("Expected int result, got %T", result)
	}
	if resultInt != 84 {
		t.Errorf("Expected 84, got %v", resultInt)
	}

	// Verify queue size
	if line.QSize() != 10 {
		t.Errorf("Expected queue size 10, got %d", line.QSize())
	}
}

func TestLine_ContextCancellation(t *testing.T) {
	var wg sync.WaitGroup

	line := NewLine(&wg, WithQSize(10), WithName("test-line-cancel"))
	line.Run()
	defer line.Stop()

	// Test with canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	callCtx := NewCallCtx(func(_ context.Context, _ any) (any, error) {
		time.Sleep(100 * time.Millisecond) // This should not complete
		return "should not reach here", nil
	}, "test")

	result, err := line.AsyncCall(ctx, callCtx)
	if err == nil {
		t.Error("Expected context cancellation error")
	}
	if result != nil {
		t.Errorf("Expected nil result on cancellation, got %v", result)
	}
}

func TestLine_MultipleAsyncCalls(t *testing.T) {
	var wg sync.WaitGroup

	line := NewLine(&wg, WithQSize(100), WithName("test-line-multiple"))
	line.Run()
	defer line.Stop()

	// Launch multiple async calls
	const numCalls = 10
	results := make(chan int, numCalls)
	errors := make(chan error, numCalls)

	for i := 0; i < numCalls; i++ {
		go func(val int) {
			ctx := context.Background()
			callCtx := NewCallCtx(func(_ context.Context, req any) (any, error) {
				input, ok := req.(int)
				if !ok {
					return nil, fmt.Errorf("expected int, got %T", req)
				}
				return input + 100, nil
			}, val)

			result, err := line.AsyncCall(ctx, callCtx)
			if err != nil {
				errors <- err
				return
			}
			resultInt, ok := result.(int)
			if !ok {
				errors <- fmt.Errorf("expected int result, got %T", result)
				return
			}
			results <- resultInt
		}(i)
	}

	// Collect results
	for i := 0; i < numCalls; i++ {
		select {
		case err := <-errors:
			t.Fatalf("Async call failed: %v", err)
		case result := <-results:
			if result < 100 || result >= 100+numCalls {
				t.Errorf("Result %d is not in expected range [100, %d)", result, 100+numCalls)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for results")
		}
	}
}
