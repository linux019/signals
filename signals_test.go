package signals_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/linux019/signals"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignal(t *testing.T) {
	testSignal := signals.NewSync[int]()

	results := make([]int, 0)
	testSignal.AddListener(func(ctx context.Context, v int) {
		results = append(results, v)
	})

	testSignal.AddListener(func(ctx context.Context, v int) {
		results = append(results, v)
	})

	ctx := context.Background()
	assert.NoError(t, testSignal.Emit(ctx, 1))
	assert.NoError(t, testSignal.Emit(ctx, 2))
	assert.NoError(t, testSignal.Emit(ctx, 3))

	require.Len(t, results, 6)
	require.Equal(t, []int{1, 1, 2, 2, 3, 3}, results)
}

func TestSignalAsync(t *testing.T) {
	var count atomic.Int32
	var wg sync.WaitGroup
	wg.Add(6)

	testSignal := signals.New[int]()
	testSignal.AddListener(func(ctx context.Context, v int) {
		time.Sleep(100 * time.Millisecond)
		count.Add(1)
		wg.Done()
	})
	testSignal.AddListener(func(ctx context.Context, v int) {
		time.Sleep(100 * time.Millisecond)
		count.Add(1)
		wg.Done()
	})

	ctx := context.Background()

	for i := 0; i < 3; i++ {
		go func(i int) {
			assert.NoError(t, testSignal.Emit(ctx, i))
		}(i)
	}

	wg.Wait()
	require.Equal(t, int32(6), count.Load())
}

// Test Async with Timeout Context. After the context is cancelled, the
// listeners should cancel their execution.
func TestSignalAsyncWithTimeout(t *testing.T) {
	var count atomic.Int32
	var timeoutCount atomic.Int32

	testSignal := signals.New[int]()
	testSignal.AddListener(func(ctx context.Context, v int) {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ctx.Done():
			timeoutCount.Add(1)
		default:
			count.Add(1)
		}
	})
	testSignal.AddListener(func(ctx context.Context, v int) {
		time.Sleep(500 * time.Millisecond)
		select {
		case <-ctx.Done():
			timeoutCount.Add(1)
		default:
			count.Add(1)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	assert.NoError(t, testSignal.Emit(ctx, 1))

	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()
	assert.NoError(t, testSignal.Emit(ctx2, 1))

	ctx3, cancel3 := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel3()
	assert.NoError(t, testSignal.Emit(ctx3, 1))

	// The code is checking if the value of the `count` variable is equal to 3 and if
	// the value of the `timeoutCount` variable is equal to 3. If either of these
	// conditions is not met, an error message is printed.
	assert.Equal(t, int32(3), count.Load())
	assert.Equal(t, int32(3), timeoutCount.Load())
}

func TestAddRemoveListener(t *testing.T) {
	testSignal := signals.New[int]()

	t.Run("AddListener", func(t *testing.T) {
		testSignal.AddListener(func(ctx context.Context, v int) {
			// Do something
		})

		testSignal.AddListener(func(ctx context.Context, v int) {
			// Do something
		}, signals.SignalType(1))

		if testSignal.Len() != 2 {
			t.Error("Count must be 2")
		}

		if count := testSignal.AddListener(func(ctx context.Context, v int) {

		}, signals.SignalType(1)); count != -1 {
			t.Error("Count must be -1")
		}
	})

	t.Run("RemoveListener", func(t *testing.T) {
		if count := testSignal.RemoveListener(signals.SignalType(1)); count != 1 {
			t.Error("Count must be 1")
		}

		if count := testSignal.RemoveListener(signals.SignalType(1)); count != -1 {
			t.Error("Count must be -1")
		}
	})

	t.Run("Reset", func(t *testing.T) {
		testSignal.Reset()
		if !testSignal.IsEmpty() {
			t.Error("Count must be 0")
		}
	})

}

// TestBaseSignal tests the BaseSignal to make sure
// Emit throws a panic because it is a base class.
func TestBaseSignal(t *testing.T) {
	testSignal := signals.BaseSignal[int]{}

	require.Error(t, testSignal.Emit(context.Background(), 1))
}
