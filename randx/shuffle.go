package randx

var (
	shuffleRandFn = RandBetweenSecure
)

// Swapper interface, each slice is a swapper, for example
//   type IntSwap []int
//   func (x IntSwap) (i, j) {
//       x[i], x[j] = x[j], x[i]
//   }
//
//   func (x IntSwap) Len() int{
//       return len(x)
//   }
type Swapper interface {
	Swap(i, j int)
	Len() int
}

// Shuffle : shuffle slice
func Shuffle(ss Swapper) {
	var size = ss.Len()
	for i := size - 1; i >= 0; i-- {
		var j = shuffleRandFn(0, i)
		ss.Swap(i, j)
	}
}

// SetShuffleRand : set shuffle rand function
func SetShuffleRand(f func(min int, max int) int) {
	shuffleRandFn = f
}
