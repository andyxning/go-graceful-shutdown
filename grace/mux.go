package grace

import "net/http"

type graceServeMux struct {
	http.ServeMux
}

func (gsm *graceServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	gsm.Handle(pattern, graceHandlerFunc(handler))
}

// NewGraceServeMux allocates and returns a new GraceServeMux.
func newGraceServeMux() *graceServeMux { return &graceServeMux{*http.NewServeMux()} }

// DefaultGraceServeMux should be register in `grace.GraceServer` with `Handler`
var DefaultGraceServeMux = newGraceServeMux()
