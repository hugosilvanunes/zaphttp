package zaphttp

import (
	"net/http"

	"go.uber.org/zap"
)

type Client struct {
	logger      *zap.Logger
	requestLog  []LogCallback
	responseLog []LogCallback
}

func newClient(logger *zap.Logger) *Client {
	return &Client{
		logger: logger,
		requestLog: []func(*Request) error{
			parseRequestBody,
		},
		responseLog: []func(*Request) error{
			parseResponseBody,
		},
	}
}

func (c *Client) AddRequestCallback(fn LogCallback) *Client {
	c.requestLog = append(c.requestLog, fn)

	return c
}

func (c *Client) AddRequestCallbacks(fns ...LogCallback) *Client {
	c.requestLog = append(c.requestLog, fns...)

	return c
}

func (c *Client) AddResponseCallback(fn LogCallback) *Client {
	c.responseLog = append(c.responseLog, fn)

	return c
}

func (c *Client) AddResponseCallbacks(fns ...LogCallback) *Client {
	c.responseLog = append(c.responseLog, fns...)

	return c
}

func (c *Client) R() *Request {
	return &Request{
		client: *c,
	}
}

func (c *Client) Middleware() func(next http.Handler) http.Handler {
	return c.R().Middleware()
}
