package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CommandArgs_with_MinimumArgumentsRequired(t *testing.T) {
	// Arrange
	const expectedValue = 3

	// Act
	sut := NewCommandArgs(MinimumArgumentsRequired(expectedValue))

	// Assert
	assert.Equal(t, expectedValue, sut.MinimumArgs)
}
