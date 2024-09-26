package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_HelloCommand_Execute(t *testing.T) {
	// Arrange
	sut := CreateHelloCommand()
	sut.SetArgs([]string{"John Doe"})

	// Act
	err := sut.Execute()

	// Assert
	assert.NoError(t, err)
}
