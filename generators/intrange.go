package generators

// IntRangeCh returns a channel that contains a sequence of ints.
// The channel is closed after the last element of the sequence
// is retrieved.
func IntRangeCh(start, count int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < count; i++ {
			ch <- start + i
		}
	}()
	return ch
}

// IntRange returns a list of integers in the specified range.
func IntRange(start, count int) []int {
	r := make([]int, count)
	for i := 0; i < count; i++ {
		r[i] = start + i
	}
	return r
}
