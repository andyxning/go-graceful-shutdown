package grace

import "net/http"

// HandleFunc registers a normal http handler function with one url path.
// This function will internally call the "DefaultServeMux.HandleFunc".
func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
