package providers

import (
	"strings"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"gopkg.in/yaml.v3"
)

type yamlProvider struct {
	data       map[string]any
	parentPath string
}

type YamlProviderOption func(*yamlProvider)

// WithParentPath sets the dotted path to a parent element in the YAML source.
func WithParentPath(path string) YamlProviderOption {
	return func(p *yamlProvider) {
		p.parentPath = path
	}
}

// NewYamlDefaultValueProvider creates a new DefaultValueProvider that reads from a YAML source.
func NewYamlDefaultValueProvider(yamlContent []byte, options ...YamlProviderOption) (types.DefaultValueProvider, error) {
	var data map[string]any
	err := yaml.Unmarshal(yamlContent, &data)
	if err != nil {
		return nil, err
	}
	provider := &yamlProvider{
		data: data,
	}
	for _, option := range options {
		option(provider)
	}
	return provider, nil
}

// GetValue retrieves the value associated with the specified key from the YAML source.
// The key is a dotted path (e.g., "service.endpoint").
// If a parent path is set, it is prepended to the key.
func (p *yamlProvider) GetValue(key string) (any, error) {
	fullKey := key
	if p.parentPath != "" {
		fullKey = p.parentPath + "." + key
	}

	parts := strings.Split(fullKey, ".")
	var current any = p.data

	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, nil
		}
		val, ok := m[part]
		if !ok {
			return nil, nil
		}
		current = val
	}

	return current, nil
}
