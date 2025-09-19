package spreadsheet

import (
	"fmt"
	"os"

	"github.com/cndrsdrmn/go-spreadsheet/internal/exitcode"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "spreadsheet",
	Short:         "Spreadsheet CLI for reading/writing",
	Long:          `A fast CLI tool built in Go for streaming rows in JSON format.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() exitcode.Code {
	if err := rootCmd.Execute(); err != nil {
		code := exitcode.FromError(err)
		fmt.Fprintln(os.Stderr, "Error:", err)
		return code
	}

	return exitcode.OK
}
