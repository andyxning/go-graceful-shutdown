package grace

import "net/http"

type HandlerFunc func(http.ResponseWriter, *http.Request)

// ServeHTTP calls the registered http handler function.
//
// All magic happends here.
//
// We track all the HTTP requests that will finally invoke the corresponding
// handler function. This is the last function that will be called before the
// actual HTTP handler function.
//
// The reason we do it here to increase and decrease the count of currently
// handling HTTP requests is that all the valid HTTP requests will be here and
// we should serve them all seriously :). The other invalid HTTP requests refer
// to those that are accepted but filtered out before reaching ServeHTTP
// function.
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defaultHTTPBarrier.Increase()
	defer defaultHTTPBarrier.Decrease()
	// f invokes the actual handler function
	f(w, r)
}
