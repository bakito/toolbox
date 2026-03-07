package http

import (
	"errors"
	"net"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCheckError(t *testing.T) {
	err := errors.New("test")
	urlErr := &url.Error{Err: err}
	wrongOpErr := &url.Error{Err: &net.OpError{Op: "foo"}}
	dialErr := &url.Error{Err: &net.OpError{Op: dialOperation}}
	tests := []struct {
		name          string
		err           error
		want          error
		wantLogFormat string
	}{
		{
			name: "should return the same error",
			err:  err,
			want: err,
		},
		{
			name: "should return the same url.Error",
			err:  urlErr,
			want: urlErr,
		},
		{
			name: "should return the url.Error if wrong OpError",
			err:  wrongOpErr,
			want: wrongOpErr,
		},
		{
			name:          "should log fatal error",
			err:           dialErr,
			want:          nil,
			wantLogFormat: msgFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logformat string
			originalLogFatalf := logFatalf
			logFatalf = func(format string, _ ...any) {
				logformat = format
			}
			defer func() { logFatalf = originalLogFatalf }()

			got := CheckError(tt.err)
			if !errors.Is(got, tt.want) {
				if diff := cmp.Diff(tt.want, got, cmpopts.EquateErrors()); diff != "" {
					t.Errorf("CheckError() mismatch (-want +got):\n%s", diff)
				}
			}

			if logformat != tt.wantLogFormat {
				t.Errorf("logformat = %v, want %v", logformat, tt.wantLogFormat)
			}
		})
	}
}
