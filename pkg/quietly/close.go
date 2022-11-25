package quietly

import "io"

func Close(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}
