package c2pa

import (
	"errors"
	"io"
)

func closeAndJoin(errp *error, closer io.Closer) {
	if closer == nil {
		return
	}
	closeErr := closer.Close()
	if closeErr == nil {
		return
	}
	if *errp == nil {
		*errp = closeErr
		return
	}
	*errp = errors.Join(*errp, closeErr)
}
