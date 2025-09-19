package csv

import (
	"encoding/csv"
	"io"
	"os"

	s "github.com/cndrsdrmn/go-spreadsheet/internal/spreadsheet"
)

type Reader struct {
	opts Options
}

// BatchStream streams rows in batches.
func (r *Reader) BatchStream(file *os.File) (<-chan s.Worksheet, error) {
	out := make(chan s.Worksheet)

	go func() {
		defer close(out)

		reader := r.open(file)

		rows := make([]s.Row, 0, r.opts.BatchSize)
		rowIndex := 0

		flush := func() {
			if len(rows) == 0 {
				return
			}

			rowStart := rowIndex - len(rows)
			rowEnd := rowIndex - 1

			batch := s.CreateBatch(rowStart, rowEnd, rows)
			sheet := s.CreateSheet(0, "Sheet 1", nil, []s.Batch{batch})

			out <- s.WrapWorksheet(sheet)

			rows = make([]s.Row, 0, r.opts.BatchSize)
		}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				out <- s.Worksheet{Err: err}
				return
			}

			rows = append(rows, s.Row{
				Sheet: 0,
				Index: rowIndex,
				Cells: record,
			})
			rowIndex++

			if len(rows) >= r.opts.BatchSize {
				flush()
			}
		}

		flush() // flush remaining rows
	}()

	return out, nil
}

// Close implements spreadsheet.Reader.
func (r *Reader) Close() error {
	return nil
}

// Options implements spreadsheet.Reader.
func (r *Reader) Options() any {
	return r.opts
}

// Read reads the entire file into memory.
func (r *Reader) Read(file *os.File) (s.Worksheet, error) {
	reader := r.open(file)

	records, err := reader.ReadAll()
	if err != nil {
		return s.Worksheet{}, err
	}

	rows := make([]s.Row, 0, len(records))
	for i, row := range records {
		rows = append(rows, s.Row{Sheet: 0, Index: i, Cells: row})
	}

	sheet := s.CreateSheet(0, "Sheet 1", rows, nil)
	return s.WrapWorksheet(sheet), nil
}

func (r *Reader) open(file *os.File) *csv.Reader {
	reader := csv.NewReader(file)
	reader.Comma = r.opts.Comma
	reader.Comment = r.opts.Comment
	reader.LazyQuotes = r.opts.LazyQuotes
	reader.TrimLeadingSpace = r.opts.TrimLeadingSpace
	return reader
}

func NewReader(opts ...Options) s.Reader {
	opt := Options{Comma: ',', BatchSize: 500}

	if len(opts) > 0 {
		opt.Merge(opts[0])
	}

	return &Reader{opts: opt}
}
