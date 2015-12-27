package grace

import "net/http"

type ServeMux struct {
	http.ServeMux
}

// HandleFunc registers a normal http handler function with one url path
func (sm *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	sm.Handle(pattern, HandlerFunc(handler))
}

// NewServeMux allocates and returns a new ServeMux instance.
func NewServeMux() *ServeMux { return &ServeMux{*http.NewServeMux()} }

// DefaultServeMux should be registered in "grace.Server" with "Handler"
var DefaultServeMux = NewServeMux()
