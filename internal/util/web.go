package util

import "net/http"

func BaseURL(r *http.Request) string {
	scheme := "http"
	if xf := r.Header.Get("Forwarded"); xf != "" {
		scheme = xf
	} else if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
