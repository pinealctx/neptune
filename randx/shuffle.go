package randx

var (
	shuffleRandFn = RandBetweenSecure
)

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
