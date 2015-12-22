package grace

import "net/http"

type graceHandlerFunc func(http.ResponseWriter, *http.Request)

// All magic happends here.
//
// We track all the HTTP requests that will finally invoke the corresponding
// handler function. This is the last function that will be called before the
// actual handler function.
//
// The reason we do it here to increase and decrease the count of currently
// handling HTTP requests is that all the valid HTTP requests will be here and
// we should serve them all seriously :). The other invalid HTTP requests refer
// to those that are accepted but filtered out before reaching
// `graceHandlerFunc.ServeHTTP`
func (f graceHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defaultGraceHTTPBarrier.Increase()
	// fmt.Print("enter\n")
	// fmt.Print(defaultGraceHTTPBarrier.GetCounter())
	defer func() {
		// fmt.Print("exit\n")
		defaultGraceHTTPBarrier.Decrease()
	}()
	// invoke the actual handler function
	f(w, r)
}
