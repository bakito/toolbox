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
	urlError := &url.Error{}
	if errors.As(err, &urlError) {
		opError := &net.OpError{}
		if errors.As(urlError.Err, &opError) {
			if opError.Op == dialOperation {
				logFatalf(msgFormat, err)
				return nil
			}
		}
	}

	return err
}
