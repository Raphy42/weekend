package http

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/palantir/stacktrace"
)

type RequestOption func(r *Request) error

func Header(key, value string) RequestOption {
	return func(r *Request) error {
		r.Headers.Set(key, value)
		return nil
	}
}

func JSONBody(value interface{}) RequestOption {
	return func(r *Request) error {
		buf, err := json.Marshal(value)
		if err != nil {
			return err
		}
		r.Body = buf
		r.ContentType = "application/json"
		return nil
	}
}

func UrlEncodedBody(raw map[string]interface{}) RequestOption {
	return func(r *Request) error {
		var values url.Values

		for key, value := range raw {
			if value == nil {
				continue
			}
			switch v := value.(type) {
			case string:
				values.Set(key, v)
			case int:
				values.Set(key, strconv.FormatInt(int64(v), 10))
			case int8:
				values.Set(key, strconv.FormatInt(int64(v), 10))
			case int16:
				values.Set(key, strconv.FormatInt(int64(v), 10))
			case int32:
				values.Set(key, strconv.FormatInt(int64(v), 10))
			case int64:
				values.Set(key, strconv.FormatInt(v, 10))
			case []string:
				values.Set(key, strings.Join(v, ","))
			case bool:
				values.Set(key, strconv.FormatBool(v))
			default:
				return stacktrace.NewError("no serialization possible for %T", v)
			}
		}
		r.Body = []byte(values.Encode())
		return nil
	}
}

func Body(value []byte, contentType string) RequestOption {
	return func(r *Request) error {
		r.Body = value
		r.ContentType = contentType
		return nil
	}
}
