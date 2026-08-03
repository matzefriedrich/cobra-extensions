package charmer

import (
	"context"
	"os"
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/providers"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/stretchr/testify/assert"
)

type envTestHandler struct {
	types.CommandName `cobra-x:"test"`
	Endpoint          string `cobra-x:"endpoint, setting-key=service.endpoint"`
	CapturedEndpoint  string
}

func (h *envTestHandler) Execute(ctx context.Context) {
	h.CapturedEndpoint = h.Endpoint
}

func Test_CommandLineApplication_uses_EnvironmentVariableDefaultValueProvider(t *testing.T) {
	// Arrange
	const envKey = "MY_APP_SERVICE_ENDPOINT"
	const envVal = "https://api.example.com"
	_ = os.Setenv(envKey, envVal)
	defer func(key string) {
		_ = os.Unsetenv(key)
	}(envKey)

	handler := &envTestHandler{}
	cmd := commands.CreateTypedCommand(handler)

	provider := providers.NewEnvironmentVariableDefaultValueProvider(providers.WithPrefix("MY_APP"))
	app := NewCommandLineApplication("app", "desc").
		WithDefaultValueProvider(provider).
		AddCommand(cmd)

	// Act
	app.root.SetArgs([]string{"test"})
	err := app.Execute(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, envVal, handler.CapturedEndpoint)
}
