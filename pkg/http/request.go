package http

import (
	"strings"
)

type Request struct {
	Headers     Headers
	Service     string
	Method      string
	Path        string
	Body        []byte
	ContentType string
}

func NewRequest(service, method, path string, opts ...RequestOption) (*Request, error) {
	request := Request{
		Service: service,
		Method:  strings.ToUpper(method),
		Path:    path,
		Body:    nil,
		Headers: NewHeaders(),
	}
	return request.apply(opts...)
}

func (r *Request) apply(opts ...RequestOption) (*Request, error) {
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}
