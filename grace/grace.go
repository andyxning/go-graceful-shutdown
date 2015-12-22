package grace

import "net/http"

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	DefaultGraceServeMux.HandleFunc(pattern, handler)
}
