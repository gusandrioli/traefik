package middlewares

import (
	"context"
	"net/http"
	"strings"
)

const (
	// StripPrefixKey is the key within the request context used to
	// store the stripped prefix
	StripPrefixKey key = "StripPrefix"
	// StripPrefixSlashKey is the key within the request context used to
	// store the stripped slash
	StripPrefixSlashKey key = "StripPrefixSlash"
	// ForwardedPrefixHeader is the default header to set prefix
	ForwardedPrefixHeader = "X-Forwarded-Prefix"
)

// StripPrefix is a middleware used to strip prefix from an URL request
type StripPrefix struct {
	Handler  http.Handler
	Prefixes []string
}

func (s *StripPrefix) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, prefix := range s.Prefixes {
		if strings.HasPrefix(r.URL.Path, prefix) {
			trailingSlash := r.URL.Path == prefix+"/"
			r.URL.Path = stripPrefix(r.URL.Path, prefix)
			if r.URL.RawPath != "" {
				r.URL.RawPath = stripPrefix(r.URL.RawPath, prefix)
			}
			s.serveRequest(w, r, strings.TrimSpace(prefix), trailingSlash)
			return
		}
	}
	http.NotFound(w, r)
}

func (s *StripPrefix) serveRequest(w http.ResponseWriter, r *http.Request, prefix string, trailingSlash bool) {
	r = r.WithContext(context.WithValue(r.Context(), StripPrefixSlashKey, trailingSlash))
	r = r.WithContext(context.WithValue(r.Context(), StripPrefixKey, prefix))
	r.Header.Add(ForwardedPrefixHeader, prefix)
	r.RequestURI = r.URL.RequestURI()
	s.Handler.ServeHTTP(w, r)
}

// SetHandler sets handler
func (s *StripPrefix) SetHandler(Handler http.Handler) {
	s.Handler = Handler
}

func stripPrefix(s, prefix string) string {
	return ensureLeadingSlash(strings.TrimPrefix(s, prefix))
}

func ensureLeadingSlash(str string) string {
	return "/" + strings.TrimPrefix(str, "/")
}
