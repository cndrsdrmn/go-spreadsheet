package spreadsheet_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"testing"

	s "github.com/cndrsdrmn/go-spreadsheet/internal/spreadsheet"
	"github.com/stretchr/testify/assert"
)

// fakeReader implements the spreadsheet.Reader interface for testing
type fakeReader struct {
	stream     bool
	withError  bool
	worksheets []s.Worksheet
	closed     bool
}

func (f *fakeReader) BatchStream(*os.File) (<-chan s.Worksheet, error) {
	if f.withError {
		return nil, errors.New("stream error")
	}
	ch := make(chan s.Worksheet, len(f.worksheets))
	for _, ws := range f.worksheets {
		ch <- ws
	}
	close(ch)
	return ch, nil
}

func (f *fakeReader) Close() error {
	f.closed = true
	return nil
}

func (f *fakeReader) Options() any {
	return map[string]string{"foo": "bar"}
}

func (f *fakeReader) Read(*os.File) (s.Worksheet, error) {
	if f.withError {
		return s.Worksheet{}, errors.New("read error")
	}
	return f.worksheets[0], nil
}

// captureStdout captures stdout output during fn execution
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	done := make(chan struct{})
	var buf bytes.Buffer
	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	fn()

	_ = w.Close()
	os.Stdout = old
	<-done
	return buf.String()
}

func TestExecuteReader_StreamMode(t *testing.T) {
	reader := &fakeReader{
		stream: true,
		worksheets: []s.Worksheet{
			{Sheets: []s.Sheet{
				{Index: 0, Name: "Sheet 1"},
				{Index: 1, Name: "Sheet 2"},
			}},
			{Sheets: []s.Sheet{
				{Index: 2, Name: "Sheet 3"},
			}},
		},
	}

	output := captureStdout(t, func() {
		err := s.ExecuteReader(reader, nil, true)
		assert.NoError(t, err)
	})

	// decode JSON lines
	dec := json.NewDecoder(bytes.NewBufferString(output))
	var count int
	for dec.More() {
		var ws s.Worksheet
		err := dec.Decode(&ws)
		assert.NoError(t, err)
		count++
	}

	assert.Equal(t, 2, count, "expected 2 worksheets")
}

func TestExecuteReader_ReadMode(t *testing.T) {
	reader := &fakeReader{
		worksheets: []s.Worksheet{
			{Sheets: []s.Sheet{
				{Index: 0, Name: "Sheet 1"},
			}},
		},
	}

	output := captureStdout(t, func() {
		err := s.ExecuteReader(reader, nil, false)
		assert.NoError(t, err)
	})

	var ws s.Worksheet
	err := json.Unmarshal([]byte(output), &ws)
	assert.NoError(t, err)
	assert.Len(t, ws.Sheets, 1)
	assert.Equal(t, "Sheet 1", ws.Sheets[0].Name)
}

func TestExecuteReader_ErrorPropagation(t *testing.T) {
	reader := &fakeReader{
		withError:  true,
		worksheets: []s.Worksheet{{}},
	}

	err := s.ExecuteReader(reader, nil, false)
	assert.EqualError(t, err, "read error")

	err = s.ExecuteReader(reader, nil, true)
	assert.EqualError(t, err, "stream error")
}
