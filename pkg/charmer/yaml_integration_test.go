package charmer

import (
	"context"
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/providers"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/stretchr/testify/assert"
)

type yamlTestHandler struct {
	types.CommandName `cobra-x:"test"`
	Servers           []string `cobra-x:"servers, setting-key=app.servers"`
	Endpoint          string   `cobra-x:"endpoint, setting-key=service.endpoint"`
	CapturedServers   []string
	CapturedEndpoint  string
}

func (h *yamlTestHandler) Execute(ctx context.Context) {
	h.CapturedServers = h.Servers
	h.CapturedEndpoint = h.Endpoint
}

func Test_CommandLineApplication_uses_YamlDefaultValueProvider(t *testing.T) {
	// Arrange
	yamlContent := []byte(`
app:
  servers:
    - server-a
    - server-b
service:
  endpoint: https://api.yaml.com
`)

	handler := &yamlTestHandler{}
	cmd := commands.CreateTypedCommand(handler)
	
	yamlProvider, err := providers.NewYamlDefaultValueProvider(yamlContent)
	assert.NoError(t, err)

	app := NewCommandLineApplication("app", "desc").
		WithDefaultValueProvider(yamlProvider).
		AddCommand(cmd)

	// Act
	app.root.SetArgs([]string{"test"})
	err = app.Execute(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"server-a", "server-b"}, handler.CapturedServers)
	assert.Equal(t, "https://api.yaml.com", handler.CapturedEndpoint)
}

func Test_CommandLineApplication_uses_YamlDefaultValueProvider_with_parent_path(t *testing.T) {
	// Arrange
	yamlContent := []byte(`
my-app:
  service:
    endpoint: https://api.yaml.com
`)

	handler := &yamlTestHandler{}
	cmd := commands.CreateTypedCommand(handler)
	
	yamlProvider, err := providers.NewYamlDefaultValueProvider(yamlContent, providers.WithParentPath("my-app"))
	assert.NoError(t, err)

	app := NewCommandLineApplication("app", "desc").
		WithDefaultValueProvider(yamlProvider).
		AddCommand(cmd)

	// Act
	app.root.SetArgs([]string{"test"})
	err = app.Execute(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "https://api.yaml.com", handler.CapturedEndpoint)
}
