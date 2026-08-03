package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_yamlProvider_GetValue_reads_nested_key(t *testing.T) {
	// Arrange
	yamlContent := []byte(`
app:
  service:
    endpoint: https://api.example.com
`)
	provider, err := NewYamlDefaultValueProvider(yamlContent)
	assert.NoError(t, err)

	// Act
	val, err := provider.GetValue("app.service.endpoint")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "https://api.example.com", val)
}

func Test_yamlProvider_GetValue_with_parent_path(t *testing.T) {
	// Arrange
	yamlContent := []byte(`
app:
  service:
    endpoint: https://api.example.com
`)
	provider, err := NewYamlDefaultValueProvider(yamlContent, WithParentPath("app"))
	assert.NoError(t, err)

	// Act
	val, err := provider.GetValue("service.endpoint")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "https://api.example.com", val)
}

func Test_yamlProvider_GetValue_reads_array(t *testing.T) {
	// Arrange
	yamlContent := []byte(`
servers:
  - s1
  - s2
  - s3
`)
	provider, err := NewYamlDefaultValueProvider(yamlContent)
	assert.NoError(t, err)

	// Act
	val, err := provider.GetValue("servers")

	// Assert
	assert.NoError(t, err)
	expected := []any{"s1", "s2", "s3"}
	assert.Equal(t, expected, val)
}

func Test_yamlProvider_GetValue_returns_nil_if_key_not_found(t *testing.T) {
	// Arrange
	yamlContent := []byte(`
app:
  service:
    endpoint: https://api.example.com
`)
	provider, err := NewYamlDefaultValueProvider(yamlContent)
	assert.NoError(t, err)

	// Act
	val, err := provider.GetValue("app.service.port")

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, val)
}
