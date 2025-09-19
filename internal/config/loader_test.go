package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cndrsdrmn/go-spreadsheet/internal/config"
	"github.com/stretchr/testify/assert"
)

type TestOptions struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func TestLoad_EmptyInput(t *testing.T) {
	var opts TestOptions

	result, err := config.Load[TestOptions]("")
	assert.NoError(t, err)
	assert.Equal(t, opts, result)
}

func TestLoad_JSONString(t *testing.T) {
	input := `{"name":"test","count":5}`

	result, err := config.Load[TestOptions](input)
	assert.NoError(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 5, result.Count)
}

func TestLoad_FileInput(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "config.json")
	content := `{"name":"file","count":10}`

	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(t, err)

	result, err := config.Load[TestOptions](filePath)
	assert.NoError(t, err)
	assert.Equal(t, "file", result.Name)
	assert.Equal(t, 10, result.Count)
}

func TestLoad_InvalidInput(t *testing.T) {
	input := "not-a-json-or-file"

	_, err := config.Load[TestOptions](input)
	assert.Error(t, err)
}
