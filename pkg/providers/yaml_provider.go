package providers

import (
	"os"
	"strings"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"gopkg.in/yaml.v3"
)

type yamlProvider struct {
	data       map[string]any
	parentPath string
}

type YamlProviderOption func(*yamlProvider) error

// WithParentPath sets the dotted path to a parent element in the YAML source.
func WithParentPath(path string) YamlProviderOption {
	return func(p *yamlProvider) error {
		p.parentPath = path
		return nil
	}
}

// RawYaml provides the YAML source as a byte array.
func RawYaml(yamlContent []byte) YamlProviderOption {
	return func(p *yamlProvider) error {
		return yaml.Unmarshal(yamlContent, &p.data)
	}
}

// YamlFile provides the YAML source from a file.
func YamlFile(filename string) YamlProviderOption {
	return func(p *yamlProvider) error {
		yamlContent, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(yamlContent, &p.data)
	}
}

// NewYamlDefaultValueProvider creates a new DefaultValueProvider that reads from a YAML source.
func NewYamlDefaultValueProvider(options ...YamlProviderOption) (types.DefaultValueProvider, error) {
	provider := &yamlProvider{
		data: make(map[string]any),
	}
	for _, option := range options {
		err := option(provider)
		if err != nil {
			return nil, err
		}
	}
	return provider, nil
}

// GetValue retrieves the value associated with the specified key from the YAML source.
// The key is a dotted path (e.g., "service.endpoint").
// If a parent path is set, it is prepended to the key.
func (p *yamlProvider) GetValue(key string) (string, error) {
	fullKey := key
	if p.parentPath != "" {
		fullKey = p.parentPath + "." + key
	}

	parts := strings.Split(fullKey, ".")
	var current any = p.data

	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return "", nil
		}
		val, ok := m[part]
		if !ok {
			return "", nil
		}
		current = val
	}

	if current == nil {
		return "", nil
	}

	switch v := current.(type) {
	case string:
		return v, nil
	default:
		// Attempt to convert to string or return empty for unsupported types (like slices)
		return "", nil
	}
}
