package http

import (
	"errors"
	"log"
	"net"
	"net/url"
)

func CheckError(err error) error {
	urlError := &url.Error{}
	if errors.As(err, &urlError) {
		opError := &net.OpError{}
		if errors.As(urlError.Err, &opError) {
			if opError.Op == "dial" {
				log.Fatalf("Network error - did you forget to set a proxy?\n%v", err)
			}
		}
	}

	return err
}
