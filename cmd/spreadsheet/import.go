package spreadsheet

import (
	"os"

	f "github.com/cndrsdrmn/go-spreadsheet/internal/factories"
	s "github.com/cndrsdrmn/go-spreadsheet/internal/spreadsheet"
	"github.com/spf13/cobra"
)

var (
	importConfig string
	importStream bool
)

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import a spreadsheet and stream JSON rows",
	Args:  cobra.ExactArgs(1),
	RunE:  runImportCommand,
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(&importConfig, "config", "c", "", "path to config file or JSON string with import options")
	importCmd.Flags().BoolVar(&importStream, "stream", false, "determine the output should be stream")
}

func runImportCommand(cmd *cobra.Command, args []string) error {
	path := args[0]

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	r, err := f.NewReader(path, importConfig)
	if err != nil {
		return err
	}
	defer r.Close()

	return s.ExecuteReader(r, file, importStream)
}
