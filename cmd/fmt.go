package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/cobra"

	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/helpers/ui"
)

type fmtCmd struct {
	meta

	print bool
	write bool
}

var (
	// fmtLong is long description of fmt command
	fmtLong = templates.LongDesc(`Formats all hcl files in a directory.
		By using the --print argument it will only print the result,
		otherwise it will by default overwrite the files.
		`)

	// fmtExample is examples for fmt command
	fmtExample = templates.Examples(`
		# Format all files in directory
		tau fmt

		# Print the result of formatting
		tau fmt --print
	`)
)

// newFmtCmd creates a new fmt command
func newFmtCmd() *cobra.Command {
	fc := &fmtCmd{}

	fmtCmd := &cobra.Command{
		Use:                   "fmt",
		Short:                 "Format all hcl files in directory",
		Long:                  fmtLong,
		Example:               fmtExample,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := fc.meta.init(args); err != nil {
				return err
			}

			return fc.run(args)
		},
	}

	f := fmtCmd.Flags()
	f.BoolVar(&fc.print, "print", false, "print the result")
	f.BoolVar(&fc.write, "write", true, "write result back to file")

	fc.addMetaFlags(fmtCmd)

	return fmtCmd
}

func (fc *fmtCmd) run(args []string) error {
	// load all sources
	files, err := fc.load()
	if err != nil {
		return err
	}

	ui.Header("Formatting files...")

	if err := files.Walk(fc.formatFile); err != nil {
		return err
	}

	ui.NewLine()

	return nil
}

func (fc *fmtCmd) formatFile(file *loader.ParsedFile) error {
	// File must be parseable as HCL native syntax before we'll try to format
	// it. If not, the formatter is likely to make drastic changes that would
	// be hard for the user to undo.
	_, syntaxDiags := hclsyntax.ParseConfig(file.Content, file.FullPath, hcl.Pos{Line: 1, Column: 1})
	if syntaxDiags.HasErrors() {
		return syntaxDiags
	}

	result := hclwrite.Format(file.Content)

	if !bytes.Equal(file.Content, result) {
		ui.Info("- %s", file.Name)

		if fc.print {
			fmt.Println(string(result))
		} else if fc.write {
			err := ioutil.WriteFile(file.FullPath, result, 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
