// Package http
package http

import (
	"errors"
	"log"
	"net"
	"net/url"
)

const (
	dialOperation = "dial"
	msgFormat     = "Network error - did you forget to set a proxy?\n%v"
)

var logFatalf = log.Fatalf

func CheckError(err error) error {
	if urlError, ok := errors.AsType[*url.Error](err); ok {
		if opError, ok := errors.AsType[*net.OpError](urlError.Err); ok {
			if opError.Op == dialOperation {
				logFatalf(msgFormat, err)
				return nil
			}
		}
	}

	return err
}
