package main

import (
	"os"

	"github.com/cndrsdrmn/go-spreadsheet/cmd/spreadsheet"
)

func main() {
	os.Exit(spreadsheet.Execute().Int())
}
