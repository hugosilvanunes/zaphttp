package zaphttp

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func parseRequestBody(r *Request) error {
	if r.RequestResult != nil {
		if err := json.Unmarshal(r.RequestRawBody, r.RequestResult); err != nil {
			return errors.Wrap(ErrParseRequestBody, err.Error())
		}
	}

	return nil
}

func parseResponseBody(r *Request) error {
	if r.ResponseResult != nil {
		if err := json.Unmarshal(r.ResponseRawBody, r.ResponseResult); err != nil {
			return errors.Wrap(ErrParseRequestBody, err.Error())
		}
	}

	return nil
}
