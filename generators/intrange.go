package generators

// IntRange returns a list of integers in the specified range.
func IntRange(start, count int) []int {
	r := make([]int, count)
	for i := 0; i < count; i++ {
		r[i] = start + i
	}
	return r
}
