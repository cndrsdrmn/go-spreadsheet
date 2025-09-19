package csv_test

import (
	"os"
	"path/filepath"
	"testing"

	csvpkg "github.com/cndrsdrmn/go-spreadsheet/formats/csv"
	"github.com/stretchr/testify/assert"
)

func writeTempCSV(t *testing.T, content string) *os.File {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), "test.csv")
	err := os.WriteFile(tmp, []byte(content), 0644)
	assert.NoError(t, err)

	f, err := os.Open(tmp)
	assert.NoError(t, err)
	return f
}

func TestReader_Read(t *testing.T) {
	file := writeTempCSV(t, "a,b,c\n1,2,3\n4,5,6\n")
	defer file.Close()

	reader := csvpkg.NewReader(csvpkg.Options{BatchSize: 2}).(*csvpkg.Reader)

	ws, err := reader.Read(file)
	assert.NoError(t, err)
	assert.Len(t, ws.Sheets, 1)

	sheet := ws.Sheets[0]
	assert.Equal(t, "Sheet 1", sheet.Name)
	assert.NotNil(t, sheet.Rows)
	assert.Len(t, *sheet.Rows, 3)
	assert.Equal(t, "a", (*sheet.Rows)[0].Cells[0])
}

func TestReader_BatchStream(t *testing.T) {
	file := writeTempCSV(t, "a,b\n1,2\n3,4\n5,6\n")
	defer file.Close()

	reader := csvpkg.NewReader(csvpkg.Options{BatchSize: 2}).(*csvpkg.Reader)

	ch, err := reader.BatchStream(file)
	assert.NoError(t, err)

	countBatches := 0
	totalRows := 0
	for ws := range ch {
		assert.NoError(t, ws.Err)
		assert.Len(t, ws.Sheets, 1)
		sheet := ws.Sheets[0]
		assert.NotNil(t, sheet.Batches)
		countBatches++
		totalRows += len((*sheet.Batches)[0].Rows)
	}

	assert.Equal(t, 2, countBatches)
	assert.Equal(t, 4, totalRows)
}

func TestReader_BatchStream_Error(t *testing.T) {
	file := writeTempCSV(t, "a,b\n1,2,3\n")
	defer file.Close()

	reader := csvpkg.NewReader(csvpkg.Options{BatchSize: 2}).(*csvpkg.Reader)

	ch, err := reader.BatchStream(file)
	assert.NoError(t, err)

	ws := <-ch
	assert.Error(t, ws.Err)
}

func TestReader_OptionsAndClose(t *testing.T) {
	reader := csvpkg.NewReader().(*csvpkg.Reader)

	opts := reader.Options().(csvpkg.Options)
	assert.Equal(t, ',', opts.Comma)
	assert.Equal(t, 500, opts.BatchSize)

	assert.NoError(t, reader.Close())
}

func TestReader_CustomDelimiter(t *testing.T) {
	file := writeTempCSV(t, "a;b\n1;2\n")
	defer file.Close()

	reader := csvpkg.NewReader(csvpkg.Options{Comma: ';', BatchSize: 2}).(*csvpkg.Reader)
	ws, err := reader.Read(file)
	assert.NoError(t, err)
	assert.Equal(t, "a", (*ws.Sheets[0].Rows)[0].Cells[0])
	assert.Equal(t, "b", (*ws.Sheets[0].Rows)[0].Cells[1])
}
