package quietly_test

import (
	"errors"
	"testing"

	"github.com/bakito/toolbox/pkg/quietly"
)

func TestClose(t *testing.T) {
	t.Run("Should close the Closer", func(t *testing.T) {
		cl := &closer{}
		quietly.Close(cl)
		if !cl.closed {
			t.Error("expected closer to be closed")
		}
	})
	t.Run("Should not fail on nil", func(*testing.T) {
		quietly.Close(nil)
	})
	t.Run("Should not fail when close return an error", func(t *testing.T) {
		cl := &closer{fail: true}
		quietly.Close(cl)
		if cl.closed {
			t.Error("expected closer to not be closed due to error")
		}
	})
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
