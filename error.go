package zaphttp

import "errors"

var (
	ErrBodyCtx          = errors.New("context body type wrong")
	ErrParseRequestBody = errors.New("context body type wrong")
)
