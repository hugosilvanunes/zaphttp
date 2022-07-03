package zaphttp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type LogCallback = func(*Request) error

type Request struct {
	ReqID              string
	Msg                string
	RawRequest         *http.Request
	RequestRawBody     []byte
	RequestResult      interface{}
	RequestMaskedBody  interface{}
	ResponseRawBody    []byte
	ResponseResult     interface{}
	ResponseMaskedBody interface{}
	client             Client
	logger             *zap.Logger
	startedAt          time.Time
	parseErr           error
}

func (rq *Request) SetRequestResult(result interface{}) *Request {
	rq.RequestResult = result

	return rq
}

func (rq *Request) SetResponseResult(result interface{}) *Request {
	rq.ResponseResult = result

	return rq
}

func (rq *Request) AddRequestCallbackR(fn LogCallback) *Request {
	rq.client.requestLog = append(rq.client.requestLog, fn)

	return rq
}

func (rq *Request) AddRequestCallbacksR(fns ...LogCallback) *Request {
	rq.client.requestLog = append(rq.client.requestLog, fns...)

	return rq
}

func (rq *Request) AddResponseCallbackR(fn LogCallback) *Request {
	rq.client.responseLog = append(rq.client.responseLog, fn)

	return rq
}

func (rq *Request) AddResponseCallbacksR(fns ...LogCallback) *Request {
	rq.client.responseLog = append(rq.client.responseLog, fns...)

	return rq
}

func (rq *Request) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			rq.startedAt = time.Now()
			rqID := middleware.GetReqID(r.Context())

			rq.logger = rq.client.logger.With(zap.String("reqID", rqID))
			rq.RawRequest = r
			rq.ReqID = rqID
			rq.Msg = fmt.Sprintf("Request: %s %s", r.Method, r.URL.Path)

			requestRawBody, err := io.ReadAll(r.Body)
			if err != nil {
				rq.logger.Debug("cannot read request body", zap.Error(err))
			}

			defer r.Body.Close()

			rq.RequestRawBody = requestRawBody

			for _, o := range rq.client.requestLog {
				if err := o(rq); err != nil {
					if errors.Is(err, ErrParseRequestBody) {
						rq.parseErr = err
					}
					rq.logger.Debug("error", zap.Error(err))
				}
			}

			logRequest(rq)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			buf := new(bytes.Buffer)
			ww.Tee(buf)

			defer func() {
				respBody, err := io.ReadAll(buf)
				if err != nil {
					rq.logger.Debug("cannot read response body", zap.Error(err))
				}

				rq.ResponseRawBody = respBody

				for _, o := range rq.client.responseLog {
					if err := o(rq); err != nil {
						rq.logger.Debug("error", zap.Error(err))
					}
				}

				logResponse(ww, rq)
			}()
			next.ServeHTTP(ww, WithBody(*rq))
		}

		return http.HandlerFunc(fn)
	}
}
