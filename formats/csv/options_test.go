package csv_test

import (
	"encoding/json"
	"testing"

	"github.com/cndrsdrmn/go-spreadsheet/formats/csv"
	"github.com/stretchr/testify/assert"
)

func TestOptions_UnmarshalJSON(t *testing.T) {
	jsonData := []byte(`{
		"batch_size": 1000,
		"delimiter": ",",
		"comment": "#",
		"lazy_quotes": true,
		"trim_leading_space": true
	}`)

	var opts csv.Options
	err := json.Unmarshal(jsonData, &opts)
	assert.NoError(t, err)
	assert.Equal(t, 1000, opts.BatchSize)
	assert.Equal(t, ',', opts.Comma)
	assert.Equal(t, '#', opts.Comment)
	assert.True(t, opts.LazyQuotes)
	assert.True(t, opts.TrimLeadingSpace)
}

func TestOptions_UnmarshalJSON_EmptyStrings(t *testing.T) {
	jsonData := []byte(`{
		"batch_size": 50,
		"delimiter": "",
		"comment": ""
	}`)

	var opts csv.Options
	err := json.Unmarshal(jsonData, &opts)
	assert.NoError(t, err)
	assert.Equal(t, 0, int(opts.Comma))
	assert.Equal(t, 0, int(opts.Comment))
	assert.Equal(t, 50, opts.BatchSize)
}

func TestOptions_Merge(t *testing.T) {
	base := csv.Options{
		BatchSize:        100,
		Comma:            ',',
		Comment:          '#',
		LazyQuotes:       false,
		TrimLeadingSpace: false,
	}

	other := csv.Options{
		BatchSize:        200,
		Comma:            ';',
		Comment:          '!',
		LazyQuotes:       true,
		TrimLeadingSpace: true,
	}

	expected := csv.Options{
		BatchSize:        200,
		Comma:            ';',
		Comment:          '!',
		LazyQuotes:       true,
		TrimLeadingSpace: true,
	}

	base.Merge(other)

	assert.EqualValues(t, expected, base)
}
