package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExtractCommandUse(t *testing.T) {

	examples := []string{"encryptCommand", "encrypt_command", "encrypt", "EncryptCommand"}

	for _, example := range examples {
		commandName := ExtractCommandUse(example)
		assert.Equal(t, "encrypt", commandName)
	}
}
