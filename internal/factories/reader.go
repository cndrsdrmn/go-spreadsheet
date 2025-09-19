package factories

import (
	"fmt"
	"path/filepath"

	"github.com/cndrsdrmn/go-spreadsheet/formats/csv"
	"github.com/cndrsdrmn/go-spreadsheet/internal/config"
	s "github.com/cndrsdrmn/go-spreadsheet/internal/spreadsheet"
)

func NewReader(path string, cfg string) (s.Reader, error) {
	switch filepath.Ext(path) {
	case ".csv":
		return createCSVReader(cfg)
	// TODO: .xlsx, .ods support
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", path)
	}
}

func createCSVReader(cfg string) (s.Reader, error) {
	opts, err := config.Load[csv.Options](cfg)
	if err != nil {
		return nil, err
	}
	return csv.NewReader(opts), nil
}
