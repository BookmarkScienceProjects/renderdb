package utils

// FillIntSlice fills all elements in the provided slice with the given val.
func FillIntSlice(slice []int, val int) {
	for i := 0; i < len(slice); i++ {
		slice[i] = val
	}
}
