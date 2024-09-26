package commands

import (
	"github.com/matzefriedrich/cobra-extensions/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DecryptCommand_Execute(t *testing.T) {

	// Arrange
	const expectedMessage = "Hello World"

	encryptedMessage, _ := utils.CaptureStdout(t, "", func() error {
		encryptCommand := CreateEncryptMessageCommand()
		encryptCommand.SetArgs([]string{"--message", expectedMessage})
		return encryptCommand.Execute()
	})

	decryptedMessage, err := utils.CaptureStdout(t, encryptedMessage, func() error {
		decryptCommand := CreateDecryptMessageCommand()
		return decryptCommand.Execute()
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, decryptedMessage)
}
