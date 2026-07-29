package reflection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_splitTagValue_splits_by_comma(t *testing.T) {
	// Arrange
	input := "part1,part2,part3"

	// Act
	result := splitTagValue(input)

	// Assert
	assert.Equal(t, []string{"part1", "part2", "part3"}, result)
}

func Test_splitTagValue_respects_single_quotes(t *testing.T) {
	// Arrange
	input := "part1, 'part2,with,comma', part3"

	// Act
	result := splitTagValue(input)

	// Assert
	assert.Equal(t, []string{"part1", " 'part2,with,comma'", " part3"}, result)
}

func Test_splitTagValue_respects_double_quotes(t *testing.T) {
	// Arrange
	input := `part1, "part2,with,comma", part3`

	// Act
	result := splitTagValue(input)

	// Assert
	assert.Equal(t, []string{"part1", ` "part2,with,comma"`, " part3"}, result)
}

func Test_splitTagValue_handles_escaped_quotes(t *testing.T) {
	// Arrange
	input := `part1, 'part2\'s', part3`

	// Act
	result := splitTagValue(input)

	// Assert
	assert.Equal(t, []string{"part1", ` 'part2\'s'`, " part3"}, result)
}

func Test_parseAttribute_returns_key_value_and_true_for_valid_input(t *testing.T) {
	// Arrange
	input := "key=value"

	// Act
	key, val, ok := parseAttribute(input)

	// Assert
	assert.True(t, ok)
	assert.Equal(t, "key", key)
	assert.Equal(t, "value", val)
}

func Test_parseAttribute_strips_single_quotes(t *testing.T) {
	// Arrange
	input := "key='value'"

	// Act
	key, val, ok := parseAttribute(input)

	// Assert
	assert.True(t, ok)
	assert.Equal(t, "key", key)
	assert.Equal(t, "value", val)
}

func Test_parseAttribute_strips_double_quotes(t *testing.T) {
	// Arrange
	input := `key="value"`

	// Act
	key, val, ok := parseAttribute(input)

	// Assert
	assert.True(t, ok)
	assert.Equal(t, "key", key)
	assert.Equal(t, "value", val)
}

func Test_parseAttribute_returns_false_for_invalid_input(t *testing.T) {
	// Arrange
	input := "invalid_input"

	// Act
	_, _, ok := parseAttribute(input)

	// Assert
	assert.False(t, ok)
}

func Test_parseCobraX_parses_full_tag(t *testing.T) {
	// Arrange
	input := "name-expr, help='Some help text', default=\"Some default\""

	// Act
	nameExpr, attrs := parseCobraX(input)

	// Assert
	assert.Equal(t, "name-expr", nameExpr)
	assert.Equal(t, "Some help text", attrs["help"])
	assert.Equal(t, "Some default", attrs["default"])
}

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
