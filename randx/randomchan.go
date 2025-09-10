package randx

import (
	"crypto/rand"
	"encoding/binary"
	"sync"
)

// RandomChan provides high-performance access to crypto/rand through channel buffering
// Uses 8-byte units for optimal performance with common operations like Uint64() and Float64()
type RandomChan struct {
	dataChan chan [8]byte  // 8-byte data channel for efficient random number generation
	stopChan chan struct{} // Stop signal for graceful shutdown
	once     sync.Once     // Ensure only close once
}

const (
	// ChanBufferSize defines channel buffer capacity: 1M * 8 bytes = 8MB total buffer
	ChanBufferSize = 1024 * 1024
	// ChanChunkSize defines how many bytes to read from crypto/rand each time: 1024 bytes
	ChanChunkSize = 1024
)

var (
	// ZeroUnit represents an all-zero 8-byte unit returned when generator is closed
	ZeroUnit = [8]byte{}
)

// NewRandomChan creates a new channel-based crypto random generator
// The generator uses a background goroutine to continuously fill the channel with random data
func NewRandomChan() *RandomChan {
	x := &RandomChan{
		dataChan: make(chan [8]byte, ChanBufferSize),
		stopChan: make(chan struct{}),
	}
	x.preload() // Preload some data into the channel
	// Start producer goroutine
	go x.producer()

	return x
}

// Close stops the random generator gracefully
// Can be called multiple times safely due to sync.Once
func (x *RandomChan) Close() {
	x.once.Do(func() {
		close(x.stopChan)
	})
}

// Float64Range returns a random float64 in [min, max) range
// Uses linear transformation to maintain uniform distribution
func (x *RandomChan) Float64Range(_min, _max float64) float64 {
	if _min >= _max {
		panic("invalid range: min must be less than max")
	}

	return _min + x.Float64()*(_max-_min)
}

// IntRange returns a random int in [min, max) range
// Convenient wrapper for common range operations
func (x *RandomChan) IntRange(_min, _max int) int {
	if _min >= _max {
		panic("invalid range: min must be less than max")
	}
	return _min + x.IntN(_max-_min)
}

// IntBetween returns a random int in [min, max] range (inclusive)
// Convenient wrapper for common range operations
func (x *RandomChan) IntBetween(_min, _max int) int {
	if _min == _max {
		return _min
	}
	if _min > _max {
		panic("invalid range: min must be less than or equal to max")
	}
	return _min + x.IntN(_max-_min+1)
}

// Float64 returns a float64 in [0, 1) range
// Uses 53 bits for proper float64 precision
func (x *RandomChan) Float64() float64 {
	// Use 53 bits for float64 precision to avoid bias
	return float64(x.Uint64()>>11) / (1 << 53)
}

// IntN returns a random int in [0, n) range
// Uses rejection sampling to ensure uniform distribution
func (x *RandomChan) IntN(n int) int {
	if n <= 0 {
		panic("invalid argument to IntN")
	}

	// Use rejection sampling for uniform distribution
	mask := uint64(1)
	for mask < uint64(n) {
		mask <<= 1
	}
	mask--

	for {
		val := x.Uint64() & mask
		if val < uint64(n) {
			return int(val)
		}
	}
}

// Uint64 reads a uint64 random number
// Most efficient method as it uses exactly one 8-byte unit
func (x *RandomChan) Uint64() uint64 {
	unit := x.readUnit()
	return binary.LittleEndian.Uint64(unit[:])
}

// ReadBytes reads n bytes from the generator
// Note: May waste some bytes for non-8-aligned sizes, but this is acceptable
// since most use cases need 8-byte multiples (Uint64, Float64, etc.)
func (x *RandomChan) ReadBytes(n int) []byte {
	if n <= 0 {
		return nil
	}

	result := make([]byte, n)
	offset := 0

	// Read in 8-byte units
	for offset < n {
		unit := x.readUnit()

		remaining := n - offset
		if remaining >= 8 {
			// Need full 8 bytes
			copy(result[offset:offset+8], unit[:])
			offset += 8
		} else {
			// Only need partial bytes (remaining bytes in unit are wasted)
			copy(result[offset:], unit[:remaining])
			offset += remaining
		}
	}

	return result
}

// readUnit reads one 8-byte unit from the channel
// Returns ZeroUnit if generator is closed and channel is empty (safer than panic)
func (x *RandomChan) readUnit() [8]byte {
	unit, ok := <-x.dataChan
	if !ok {
		// Channel is closed, return zero unit
		return ZeroUnit
	}
	return unit
}

// preload producer
func (x *RandomChan) preload() {
	chunk := make([]byte, ChanChunkSize)
	for {
		n, err := rand.Read(chunk)
		if err != nil {
			// Skip this round on error and try again
			continue
		}

		for i := 0; i+8 <= n; i += 8 {
			var unit [8]byte
			copy(unit[:], chunk[i:i+8])

			select {
			case x.dataChan <- unit:
				// Successfully sent
			default:
				return
			}
		}
	}
}

// producer continuously reads from crypto/rand and writes to channel
// Runs in background goroutine and blocks naturally when channel is full
func (x *RandomChan) producer() {
	defer close(x.dataChan)

	chunk := make([]byte, ChanChunkSize)

	for {
		select {
		case <-x.stopChan:
			return
		default:
			// Read from crypto/rand
			n, err := rand.Read(chunk)
			if err != nil {
				// Skip this round on error and try again
				continue
			}

			// Split into 8-byte units and send to channel
			// Channel will naturally block when full, providing backpressure
			for i := 0; i+8 <= n; i += 8 {
				var unit [8]byte
				copy(unit[:], chunk[i:i+8])

				select {
				case x.dataChan <- unit:
					// Successfully sent
				case <-x.stopChan:
					return
				}
			}
		}
	}
}
