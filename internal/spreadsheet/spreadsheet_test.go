package spreadsheet_test

import (
	"testing"

	s "github.com/cndrsdrmn/go-spreadsheet/internal/spreadsheet"
	"github.com/stretchr/testify/assert"
)

func TestCreateBatch(t *testing.T) {
	rows := []s.Row{
		{Sheet: 0, Index: 1, Cells: []string{"A1", "B1"}},
		{Sheet: 0, Index: 2, Cells: []string{"A2", "B2"}},
	}

	batch := s.CreateBatch(1, 2, rows)

	assert.Equal(t, 1, batch.RowStart)
	assert.Equal(t, 2, batch.RowEnd)
	assert.Equal(t, rows, batch.Rows)
}

func TestCreateSheet_WithRowsAndBatches(t *testing.T) {
	rows := []s.Row{{Sheet: 0, Index: 1, Cells: []string{"A1"}}}
	batches := []s.Batch{{RowStart: 1, RowEnd: 1, Rows: rows}}

	sheet := s.CreateSheet(0, "Sheet1", rows, batches)

	assert.Equal(t, 0, sheet.Index)
	assert.Equal(t, "Sheet1", sheet.Name)
	assert.NotNil(t, sheet.Rows)
	assert.Equal(t, rows, *sheet.Rows)
	assert.NotNil(t, sheet.Batches)
	assert.Equal(t, batches, *sheet.Batches)
}

func TestCreateSheet_WithoutRowsAndBatches(t *testing.T) {
	sheet := s.CreateSheet(1, "EmptySheet", nil, nil)

	assert.Nil(t, sheet.Rows)
	assert.Nil(t, sheet.Batches)
}

func TestWrapWorksheet(t *testing.T) {
	sheet1 := s.CreateSheet(0, "Sheet1", nil, nil)
	sheet2 := s.CreateSheet(1, "Sheet2", nil, nil)

	ws := s.WrapWorksheet(sheet1, sheet2)

	assert.Len(t, ws.Sheets, 2)
	assert.Equal(t, "Sheet1", ws.Sheets[0].Name)
	assert.Equal(t, "Sheet2", ws.Sheets[1].Name)
	assert.Nil(t, ws.Err)
}
