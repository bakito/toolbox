package quietly_test

import (
	"errors"
	"testing"

	"github.com/bakito/toolbox/pkg/quietly"
)

func TestClose(t *testing.T) {
	tests := []struct {
		name           string
		closer         *closer
		shouldBeClosed bool
	}{
		{
			name:           "Should close the Closer",
			closer:         &closer{},
			shouldBeClosed: true,
		},
		{
			name:           "Should not fail on nil",
			closer:         nil,
			shouldBeClosed: false,
		},
		{
			name:           "Should not fail when close return an error",
			closer:         &closer{fail: true},
			shouldBeClosed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.closer == nil {
				quietly.Close(nil)
			} else {
				quietly.Close(tt.closer)
				if tt.closer.closed != tt.shouldBeClosed {
					t.Errorf("expected closed to be %v, got %v", tt.shouldBeClosed, tt.closer.closed)
				}
			}
		})
	}
}

type closer struct {
	closed bool
	fail   bool
}

func (c *closer) Close() error {
	if c.fail {
		return errors.New("failed")
	}
	c.closed = true
	return nil
}
