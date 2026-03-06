package http

import (
	"errors"
	"net"
	"net/url"
	"testing"
)

func TestCheckError(t *testing.T) {
	t.Run("should return the same error", func(t *testing.T) {
		err := errors.New("test")
		err2 := CheckError(err)
		if !errors.Is(err2, err) {
			t.Errorf("expected %v, got %v", err, err2)
		}
	})
	t.Run("should return the same url.Error", func(t *testing.T) {
		err := errors.New("test")
		urlErr := &url.Error{Err: err}
		err2 := CheckError(urlErr)
		if !errors.Is(err2, urlErr) {
			t.Errorf("expected %v, got %v", urlErr, err2)
		}
	})
	t.Run("should return the url.Error if wrong OpError", func(t *testing.T) {
		urlErr := &url.Error{Err: &net.OpError{Op: "foo"}}
		err2 := CheckError(urlErr)
		if !errors.Is(err2, urlErr) {
			t.Errorf("expected %v, got %v", urlErr, err2)
		}
	})
	t.Run("should log fatal error", func(t *testing.T) {
		var logformat string
		originalLogFatalf := logFatalf
		defer func() { logFatalf = originalLogFatalf }()
		logFatalf = func(format string, _ ...any) {
			logformat = format
		}
		urlErr := &url.Error{Err: &net.OpError{Op: dialOperation}}
		err2 := CheckError(urlErr)
		if err2 != nil {
			t.Errorf("expected nil error, got %v", err2)
		}
		if logformat != msgFormat {
			t.Errorf("expected log format %q, got %q", msgFormat, logformat)
		}
	})
}
