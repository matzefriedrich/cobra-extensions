package charmer

import (
	"context"
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/stretchr/testify/assert"
)

type mockProvider struct {
	values map[string]string
}

func (m *mockProvider) GetValue(key string) (*string, error) {
	if val, ok := m.values[key]; ok {
		return &val, nil
	}
	return nil, nil
}

type testHandler struct {
	types.BaseCommand `cobra-x:"test"`
	MyFlag            string `cobra-x:"my-flag, setting-key=my-setting"`
	CapturedValue     string
}

func (h *testHandler) Execute(_ context.Context) {
	h.CapturedValue = h.MyFlag
}

func Test_CommandLineApplication_injects_DefaultValueProvider_to_handler(t *testing.T) {
	// Arrange
	provider := &mockProvider{
		values: map[string]string{
			"my-setting": "provider-value",
		},
	}
	handler := &testHandler{}
	cmd := commands.CreateTypedCommand(handler)

	app := NewCommandLineApplication("app", "app desc").
		WithDefaultValueProvider(provider).
		AddCommand(cmd)

	// Act
	app.root.SetArgs([]string{"test"})
	err := app.Execute(t.Context())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "provider-value", handler.CapturedValue)
}

func Test_CommandLineApplication_prefers_CLI_flag_over_provider(t *testing.T) {
	// Arrange
	provider := &mockProvider{
		values: map[string]string{
			"my-setting": "provider-value",
		},
	}
	handler := &testHandler{}
	cmd := commands.CreateTypedCommand(handler)

	app := NewCommandLineApplication("app", "app desc").
		WithDefaultValueProvider(provider).
		AddCommand(cmd)

	// Act
	app.root.SetArgs([]string{"test", "--my-flag", "cli-value"})
	err := app.Execute(t.Context())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "cli-value", handler.CapturedValue)
}
