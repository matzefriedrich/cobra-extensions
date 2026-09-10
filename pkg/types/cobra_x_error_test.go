package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CobraXError_Error_returns_the_error_message(t *testing.T) {
	// Arrange
	sut := NewCobraXError("boom")

	// Act
	actual := sut.Error()

	// Assert
	assert.Equal(t, "boom", actual)
}

func Test_CobraXError_Unwrap_returns_nil_when_no_cause_is_set(t *testing.T) {
	// Arrange
	sut := NewCobraXError("boom")

	// Act
	actual := sut.Unwrap()

	// Assert
	assert.Nil(t, actual)
}

func Test_CobraXError_Unwrap_returns_the_underlying_cause(t *testing.T) {
	// Arrange
	cause := errors.New("root cause")
	sut := NewCobraXError("boom", WithCause(cause))

	// Act
	actual := sut.Unwrap()

	// Assert
	assert.Equal(t, cause, actual)
}

func Test_CobraXError_Is_returns_true_when_the_messages_match(t *testing.T) {
	// Arrange
	sut := NewCobraXError("boom")
	target := NewCobraXError("boom")

	// Act
	actual := sut.Is(target)

	// Assert
	assert.True(t, actual)
}

func Test_CobraXError_Is_returns_false_when_the_messages_differ(t *testing.T) {
	// Arrange
	sut := NewCobraXError("boom")
	target := NewCobraXError("bang")

	// Act
	actual := sut.Is(target)

	// Assert
	assert.False(t, actual)
}

func Test_NewCobraXError_creates_an_error_with_the_given_message(t *testing.T) {
	// Arrange
	const expectedMessage = "boom"

	// Act
	sut := NewCobraXError(expectedMessage)

	// Assert
	assert.Equal(t, expectedMessage, sut.msg)
}

func Test_NewCobraXError_applies_options_to_the_created_error(t *testing.T) {
	// Arrange
	cause := errors.New("root cause")

	// Act
	sut := NewCobraXError("boom", WithCause(cause))

	// Assert
	assert.Equal(t, cause, sut.cause)
}

func Test_WithCause_sets_the_cause_of_the_error(t *testing.T) {
	// Arrange
	cause := errors.New("root cause")
	target := NewCobraXError("boom")

	// Act
	WithCause(cause)(target)

	// Assert
	assert.Equal(t, cause, target.cause)
}
