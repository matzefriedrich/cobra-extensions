package reflection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseFlagNameExpression_parses_long_and_short_names(t *testing.T) {
	// Arrange
	input := "--name|-n"

	// Act
	name, shorthand := parseFlagNameExpression(input)

	// Assert
	assert.Equal(t, "name", name)
	assert.Equal(t, "n", shorthand)
}

func Test_parseFlagNameExpression_handles_only_long_name(t *testing.T) {
	// Arrange
	input := "--name"

	// Act
	name, shorthand := parseFlagNameExpression(input)

	// Assert
	assert.Equal(t, "name", name)
	assert.Equal(t, "", shorthand)
}

func Test_parseFlagNameExpression_handles_only_short_name(t *testing.T) {
	// Arrange
	input := "-n"

	// Act
	name, shorthand := parseFlagNameExpression(input)

	// Assert
	assert.Equal(t, "n", name)
	assert.Equal(t, "n", shorthand)
}

func Test_parseFlagNameExpression_handles_simple_name_without_dashes(t *testing.T) {
	// Arrange
	input := "name"

	// Act
	name, shorthand := parseFlagNameExpression(input)

	// Assert
	assert.Equal(t, "name", name)
	assert.Equal(t, "", shorthand)
}
