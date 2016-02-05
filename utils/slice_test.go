package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils_FillIntSlice_BigSlice_SetsAllElements(t *testing.T) {
	// Arrange
	slice := make([]int, 10001)

	// Act
	FillIntSlice(slice, 1337)

	// Assert
	for _, v := range slice {
		assert.Equal(t, v, 1337)
	}
}
