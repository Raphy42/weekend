package http

import "strings"

type Headers map[string]string

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h *Headers) Set(key, value string) {
	(*h)[strings.ToLower(key)] = value
}

func (h *Headers) Get(key string) (string, bool) {
	v, ok := (*h)[strings.ToLower(key)]
	return v, ok
}
