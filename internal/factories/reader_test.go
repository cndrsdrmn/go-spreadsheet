package factories_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cndrsdrmn/go-spreadsheet/formats/csv"
	f "github.com/cndrsdrmn/go-spreadsheet/internal/factories"
	"github.com/stretchr/testify/assert"
)

func TestNewReader_CSV_EmptyConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.csv")
	os.WriteFile(path, []byte("a,b,c\n1,2,3"), 0644)

	r, err := f.NewReader(path, "")
	assert.NoError(t, err)
	assert.NotNil(t, r)

	_, ok := r.(*csv.Reader)
	assert.True(t, ok, "expected CSV reader type")
}

func TestNewReader_CSV_JSONConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.csv")
	os.WriteFile(path, []byte("a,b,c\n1,2,3"), 0644)

	cfg := `{"batch_size":2,"delimiter":",","lazy_quotes":true}`

	r, err := f.NewReader(path, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, r)

	reader, ok := r.(*csv.Reader)
	assert.True(t, ok, "expected CSV reader type")

	opts, ok := reader.Options().(csv.Options)
	assert.True(t, ok, "expected CSV Options type")

	assert.Equal(t, 2, opts.BatchSize)
	assert.Equal(t, ',', opts.Comma)
	assert.True(t, opts.LazyQuotes)
}

func TestNewReader_UnsupportedExtension(t *testing.T) {
	path := "file.txt"

	r, err := f.NewReader(path, "")
	assert.Nil(t, r)
	assert.Error(t, err)
}
