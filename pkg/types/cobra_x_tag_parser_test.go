package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_compactTagParser_ParseField_parses_name_expression(t *testing.T) {
	// Arrange
	field := reflect.StructField{Name: "Name", Tag: `cobra-x:"-n|--name"`}
	sut := NewCompactTagParser()

	// Act
	name, attributes := sut.ParseField(field)

	// Assert
	assert.Equal(t, "-n|--name", name)
	assert.Empty(t, attributes)
}

func Test_compactTagParser_ParseField_preserves_quoted_attribute_values(t *testing.T) {
	// Arrange
	field := reflect.StructField{Name: "Name", Tag: `cobra-x:"name, help='The name to greet', default='World'"`}
	sut := NewCompactTagParser()

	// Act
	name, attributes := sut.ParseField(field)

	// Assert
	assert.Equal(t, "name", name)
	assert.Equal(t, "The name to greet", attributes["help"])
	assert.Equal(t, "World", attributes["default"])
}

func Test_standardTagParser_ParseField_parses_name_and_shorthand_keys(t *testing.T) {
	// Arrange
	field := reflect.StructField{Name: "Name", Tag: `cobra-x:"--name" cobra-x-shorthand:"n"`}
	sut := NewStandardTagParser()

	// Act
	name, attributes := sut.ParseField(field)

	// Assert
	assert.Equal(t, "-n|--name", name)
	assert.Empty(t, attributes)
}

func Test_standardTagParser_ParseField_parses_flag_attribute_keys(t *testing.T) {
	// Arrange
	field := reflect.StructField{Name: "Name", Tag: `cobra-x:"--name" cobra-x-usage:"The name to greet" cobra-x-default:"World" cobra-x-setting-key:"name-setting"`}
	sut := NewStandardTagParser()

	// Act
	name, attributes := sut.ParseField(field)

	// Assert
	assert.Equal(t, "--name", name)
	assert.Equal(t, "The name to greet", attributes["usage"])
	assert.Equal(t, "World", attributes["default"])
	assert.Equal(t, "name-setting", attributes["setting-key"])
}

func Test_standardTagParser_ParseField_parses_command_help_and_description_keys(t *testing.T) {
	// Arrange
	field := reflect.StructField{Name: "Name", Tag: `cobra-x:"greet" cobra-x-help:"Greets someone" cobra-x-description:"Greets the given person warmly"`}
	sut := NewStandardTagParser()

	// Act
	name, attributes := sut.ParseField(field)

	// Assert
	assert.Equal(t, "greet", name)
	assert.Equal(t, "Greets someone", attributes["help"])
	assert.Equal(t, "Greets the given person warmly", attributes["description"])
}

func Test_standardTagParser_ParseField_returns_empty_for_missing_cobra_x_tag(t *testing.T) {
	// Arrange
	field := reflect.StructField{Name: "Name", Tag: `json:"name"`}
	sut := NewStandardTagParser()

	// Act
	name, attributes := sut.ParseField(field)

	// Assert
	assert.Empty(t, name)
	assert.Empty(t, attributes)
}
