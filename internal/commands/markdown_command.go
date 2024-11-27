package commands

import (
	"github.com/matzefriedrich/cobra-extensions/internal/utils"
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"log"
	"os"
	"path/filepath"
)

type markdownDocsCommand struct {
	use              types.CommandName `flag:"markdown" short:"Exports Markdown documentation to the specified folder" description:"Exports Markdown documentation to the specified folder"`
	OutputFolderPath string            `flag:"output" usage:"The output folder for markdown documentation" default:"."`
	root             *cobra.Command
}

func (m *markdownDocsCommand) Execute() {

	outputPath, _ := filepath.Abs(m.OutputFolderPath)

	ext := filepath.Ext(outputPath)
	switch ext {
	case ".md", ".markdown":
		m.generateSingleMarkdownFile(outputPath)
	default:
		err := doc.GenMarkdownTree(m.root, outputPath)
		if err != nil {
			log.Fatalf("Error generating markdown documentation: %v", err)
		}
	}

}

func (m *markdownDocsCommand) generateSingleMarkdownFile(outputPath string) {

	writer, _ := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer func(writer *os.File) {
		_ = writer.Close()
	}(writer)

	stack := utils.MakeStack[*cobra.Command]()
	stack.Push(m.root)

	for stack.Any() {
		next := stack.Pop()
		for _, c := range next.Commands() {
			if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
				continue
			}
			stack.Push(c)
		}

		_ = doc.GenMarkdownCustom(next, writer, func(s string) string {
			return s
		})

		_, _ = writer.WriteString("\n")
	}

	_ = writer.Sync()
}

var _ types.TypedCommand = (*markdownDocsCommand)(nil)

// NewMarkdownDocsCommand creates a new Cobra command for exporting Markdown documentation to a specified folder.
func NewMarkdownDocsCommand(root *cobra.Command) *cobra.Command {
	instance := &markdownDocsCommand{
		root: root,
	}
	return commands.CreateTypedCommand(instance)
}
