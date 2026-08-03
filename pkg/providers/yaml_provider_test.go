package providers

import (
	"os"
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
	provider, err := NewYamlDefaultValueProvider(RawYaml(yamlContent))
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
	provider, err := NewYamlDefaultValueProvider(RawYaml(yamlContent), WithParentPath("app"))
	assert.NoError(t, err)

	// Act
	val, err := provider.GetValue("service.endpoint")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "https://api.example.com", val)
}

func Test_yamlProvider_GetValue_returns_nil_if_key_not_found(t *testing.T) {
	// Arrange
	yamlContent := []byte(`
app:
  service:
    endpoint: https://api.example.com
`)
	provider, err := NewYamlDefaultValueProvider(RawYaml(yamlContent))
	assert.NoError(t, err)

	// Act
	val, err := provider.GetValue("app.service.port")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "", val)
}

func Test_yamlProvider_GetValue_reads_from_file(t *testing.T) {
	// Arrange
	yamlContent := []byte(`
app:
  service:
    endpoint: https://api.file.com
`)
	tmpFile, err := os.CreateTemp("", "test*.yaml")
	assert.NoError(t, err)
	defer func(name string) {
		_ = os.Remove(name)
	}(tmpFile.Name())

	_, err = tmpFile.Write(yamlContent)
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)

	provider, err := NewYamlDefaultValueProvider(YamlFile(tmpFile.Name()))
	assert.NoError(t, err)

	// Act
	val, err := provider.GetValue("app.service.endpoint")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "https://api.file.com", val)
}
