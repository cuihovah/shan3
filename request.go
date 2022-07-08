package shan3

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func QueryParse(r *http.Request) map[string]string {
	ret := make(map[string]string)
	data := strings.Split(r.URL.RawQuery, "&")
	for _, v := range data {
		kv := strings.Split(v, "=")
		if len(kv) == 1 {
			ret[kv[0]] = "true"
		}
		if len(kv) == 2 {
			ret[kv[0]], _ = url.QueryUnescape(kv[1])
		}
	}
	return ret
}

func BodyBuffer(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	return body, err
}

func GetMethodName(r *http.Request) string {
	data := r.URL.Query()
	for key, value := range data {
		if key == "method" {
			return value[0]
		}
	}
	return ""
}

func WithValue(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	ctx = context.WithValue(ctx, "request", r)
	ctx = context.WithValue(ctx, "response", w)
	return ctx
}
