package ternary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValueOrDefault_returns_the_value_when_the_condition_is_true(t *testing.T) {
	// Arrange
	value := 42
	isPositive := func(candidate int) bool {
		return candidate > 0
	}

	// Act
	actual := ValueOrDefault(value, isPositive, 0)

	// Assert
	assert.Equal(t, 42, actual)
}

func Test_ValueOrDefault_returns_the_false_value_when_the_condition_is_false(t *testing.T) {
	// Arrange
	value := -1
	isPositive := func(candidate int) bool {
		return candidate > 0
	}

	// Act
	actual := ValueOrDefault(value, isPositive, 0)

	// Assert
	assert.Equal(t, 0, actual)
}

func Test_NotNilOrWhitespace_returns_true_for_a_non_blank_string(t *testing.T) {
	// Arrange
	value := "encrypt"

	// Act
	actual := NotNilOrWhitespace(value)

	// Assert
	assert.True(t, actual)
}

func Test_NotNilOrWhitespace_returns_false_for_an_empty_string(t *testing.T) {
	// Arrange
	value := ""

	// Act
	actual := NotNilOrWhitespace(value)

	// Assert
	assert.False(t, actual)
}

func Test_NotNilOrWhitespace_returns_false_for_a_whitespace_only_string(t *testing.T) {
	// Arrange
	value := " \t "

	// Act
	actual := NotNilOrWhitespace(value)

	// Assert
	assert.False(t, actual)
}
