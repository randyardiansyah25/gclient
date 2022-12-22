package err

import "errors"

var (
	ErrRefused = errors.New("connection refuse")
	ErrTimeout = errors.New("connection timeout")
)
