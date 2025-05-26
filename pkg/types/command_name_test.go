package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CommandName_String_returns_the_name_field(t *testing.T) {
	// Arrange
	sut := CommandName{
		name: "test",
	}

	// Act
	actual := sut.String()

	// Assert
	assert.Equal(t, "test", actual)
}
