package providers

import (
	"os"
	"strings"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
)

type envProvider struct {
	prefix string
}

type EnvProviderOption func(*envProvider)

// WithPrefix sets the prefix for environment variables.
func WithPrefix(prefix string) EnvProviderOption {
	return func(p *envProvider) {
		p.prefix = prefix
	}
}

// NewEnvironmentVariableDefaultValueProvider creates a new DefaultValueProvider that reads from environment variables.
func NewEnvironmentVariableDefaultValueProvider(options ...EnvProviderOption) types.DefaultValueProvider {
	provider := &envProvider{}
	for _, option := range options {
		option(provider)
	}
	return provider
}

// GetValue retrieves the value associated with the specified key from environment variables.
// The key is transformed by replacing "." with "_" and converting it to uppercase.
// If a prefix is set, it is prepended to the key with an underscore.
func (p *envProvider) GetValue(key string) (string, error) {
	envKey := strings.ReplaceAll(key, ".", "_")
	envKey = strings.ToUpper(envKey)

	if p.prefix != "" {
		prefix := strings.TrimSuffix(p.prefix, "_")
		envKey = prefix + "_" + envKey
	}

	val, ok := os.LookupEnv(envKey)
	if !ok {
		return "", nil
	}

	return val, nil
}
