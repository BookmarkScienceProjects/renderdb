package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstError_Empty_ReturnsNil(t *testing.T) {
	assert.NoError(t, FirstError())
}

func TestFirstError_NoErrors_ReturnsNil(t *testing.T) {
	assert.NoError(t, FirstError(nil, nil, nil))
}

func TestFirstError_TwoErrors_ReturnsFirst(t *testing.T) {
	assert.Equal(t, errors.New("1"), FirstError(nil, errors.New("1"), errors.New("2")))
}
