package zaphttp

import (
	"fmt"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func newLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		logger = zap.NewNop()
	}

	return logger
}

func logRequest(r *Request) {
	fields := make([]zap.Field, 0)

	fields = append(fields,
		zap.String("from", "client"),
		zap.String("to", "gateway"),
		zap.String("method", r.RawRequest.Method),
		zap.String("path", r.RawRequest.URL.Path),
		zap.String("protocol", r.RawRequest.Proto),
	)

	if len(r.RawRequest.Header) > 0 {
		fields = append(fields, zap.Any("header", headerLogField(r.RawRequest.Header)))
	}

	if r.RequestMaskedBody != nil {
		fields = append(fields, zap.Any("body", r.RequestMaskedBody))
	} else {
		fields = append(fields, zap.ByteString("body", r.RequestRawBody))
	}

	r.logger.Info(r.Msg, fields...)
}

func logResponse(ww middleware.WrapResponseWriter, rq *Request) {
	msg := fmt.Sprintf("Response: %d", ww.Status())

	if rq.Msg != "" {
		msg = fmt.Sprintf("%s - %s", msg, rq.Msg)
	}

	fieldsResponse := make([]zap.Field, 0)
	fieldsResponse = append(fieldsResponse,
		zap.String("from", "gateway"),
		zap.String("to", "client"),
		zap.Int("status", ww.Status()),
		zap.Float64("elapsed", time.Since(rq.startedAt).Seconds()),
	)

	if rq.ResponseMaskedBody != nil {
		fieldsResponse = append(fieldsResponse, zap.Any("body", rq.ResponseMaskedBody))
	} else {
		fieldsResponse = append(fieldsResponse, zap.ByteString("body", rq.ResponseRawBody))
	}

	if len(ww.Header()) > 0 {
		fieldsResponse = append(fieldsResponse, zap.Any("header", headerLogField(ww.Header())))
	}

	rq.logger.Info(msg, fieldsResponse...)
}
