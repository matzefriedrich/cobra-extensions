package providers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_envProvider_GetValue_transforms_key_without_prefix(t *testing.T) {
	// Arrange
	key := "service.endpoint"
	expectedEnvKey := "SERVICE_ENDPOINT"
	expectedValue := "http://localhost:8080"

	_ = os.Setenv(expectedEnvKey, expectedValue)
	defer func(key string) {
		_ = os.Unsetenv(key)
	}(expectedEnvKey)

	provider := NewEnvironmentVariableDefaultValueProvider()

	// Act
	val, err := provider.GetValue(key)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, val)
}

func Test_envProvider_GetValue_transforms_key_with_prefix(t *testing.T) {
	// Arrange
	key := "service.endpoint"
	prefix := "MY_APP"
	expectedEnvKey := "MY_APP_SERVICE_ENDPOINT"
	expectedValue := "https://api.openai.com"

	_ = os.Setenv(expectedEnvKey, expectedValue)
	defer func(key string) {
		_ = os.Unsetenv(key)
	}(expectedEnvKey)

	provider := NewEnvironmentVariableDefaultValueProvider(WithPrefix(prefix))

	// Act
	val, err := provider.GetValue(key)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, val)
}

func Test_envProvider_GetValue_returns_nil_if_env_var_is_not_set(t *testing.T) {
	// Arrange
	key := "non.existent"
	provider := NewEnvironmentVariableDefaultValueProvider()

	// Act
	val, err := provider.GetValue(key)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "", val)
}

func Test_envProvider_GetValue_handles_prefix_with_trailing_underscore(t *testing.T) {
	// Arrange
	key := "setting"
	prefix := "APP_"
	expectedEnvKey := "APP_SETTING"
	expectedValue := "val"

	_ = os.Setenv(expectedEnvKey, expectedValue)
	defer func(key string) {
		_ = os.Unsetenv(key)
	}(expectedEnvKey)

	provider := NewEnvironmentVariableDefaultValueProvider(WithPrefix(prefix))

	// Act
	val, err := provider.GetValue(key)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, val)
}
