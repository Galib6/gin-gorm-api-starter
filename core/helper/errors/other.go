package errs

import "errors"

var (
	ErrInvalidPage    = errors.New("entered page is invalid")
	ErrInvalidSort    = errors.New("entered sort parameter is invalid")
	ErrInvalidInclude = errors.New("entered include parameter is invalid")
)
