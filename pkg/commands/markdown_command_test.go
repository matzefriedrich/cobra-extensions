package commands

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func Test_MarkdownCommand_Execute_renders_multiple_files_if_folder_path_specified(t *testing.T) {

	// Arrange
	app := &cobra.Command{}
	sut := NewMarkdownDocsCommand(app)

	app.AddCommand(sut)
	app.SetArgs([]string{"markdown", "--output", os.TempDir()})

	// Act
	err := app.Execute()

	// Assert
	assert.NoError(t, err)
}

func Test_MarkdownCommand_Execute_renders_single_file_if_markdown_file_path_specified(t *testing.T) {

	// Arrange
	app := &cobra.Command{}
	sut := NewMarkdownDocsCommand(app)

	app.AddCommand(sut)
	app.SetArgs([]string{"markdown", "--output", filepath.Join(os.TempDir(), "test.md")})

	// Act
	err := app.Execute()

	// Assert
	assert.NoError(t, err)
}
