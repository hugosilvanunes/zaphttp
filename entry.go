package zaphttp

import (
	"context"
	"net/http"
)

type contextKey struct {
	name string
}

var BodyEntryCtxKey = &contextKey{"BodyEntry"}

type BodyEntry struct {
	raw    []byte
	parsed interface{}
}

func NewBodyEntry(raw []byte, parsed interface{}) *BodyEntry {
	return &BodyEntry{
		raw:    raw,
		parsed: parsed,
	}
}

func Body(ctx context.Context) interface{} {
	entry := getBodyEntryCtx(ctx)
	if entry == nil {
		return nil
	}

	return entry.parsed
}

func RawBody(ctx context.Context) []byte {
	entry := getBodyEntryCtx(ctx)
	if entry == nil {
		return nil
	}

	return entry.raw
}

func WithBody(rq Request) *http.Request {
	r := rq.RawRequest

	bodyEntry := NewBodyEntry(rq.RequestRawBody, rq.RequestResult)

	r = r.WithContext(context.WithValue(r.Context(), BodyEntryCtxKey, *bodyEntry))

	return r
}

func getBodyEntryCtx(ctx context.Context) *BodyEntry {
	entry, _ := ctx.Value(BodyEntryCtxKey).(BodyEntry)

	return &entry
}
