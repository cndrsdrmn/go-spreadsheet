package exitcode_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	e "github.com/cndrsdrmn/go-spreadsheet/internal/exitcode"
	"github.com/stretchr/testify/assert"
)

func TestCode_Int(t *testing.T) {
	assert.Equal(t, 0, e.OK.Int())
	assert.Equal(t, 1, e.Fail.Int())
	assert.Equal(t, 2, e.InvalidArgs.Int())
	assert.Equal(t, 3, e.IOError.Int())
	assert.Equal(t, 4, e.NotFound.Int())
}

func TestFromError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want e.Code
	}{
		{"nil error", nil, e.OK},
		{"not exist", os.ErrNotExist, e.NotFound},
		{"permission", os.ErrPermission, e.IOError},
		{"other error", errors.New("some error"), e.Fail},
		{"wrapped not exist", fmt.Errorf("wrapped: %w", os.ErrNotExist), e.NotFound},
		{"wrapped permission", fmt.Errorf("wrapped: %w", os.ErrPermission), e.IOError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := e.FromError(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
