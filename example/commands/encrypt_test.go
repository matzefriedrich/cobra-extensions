package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CryptCommand_Execute(t *testing.T) {

	// Arrange
	sut := CreateEncryptMessageCommand()
	sut.SetArgs([]string{"--message", "Hello World"})
	// Act
	err := sut.Execute()

	// Assert
	assert.NoError(t, err)
}
