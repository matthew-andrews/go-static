package staticresponsewriter

import (
	"log"
	"net/http"
)

// Define a ‘StaticResponseWriter’ in order do some manipulation and logging
// of the response immediately before the header is flushed in ‘WriteHeader’
// otherwise, this is just straight up inheritance
type StaticResponseWriter struct {
	Headers        map[string]string
	ResponseWriter http.ResponseWriter
	Path           string
}

func (w StaticResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w StaticResponseWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func (w StaticResponseWriter) WriteHeader(status int) {
	for header, value := range w.Headers {
		w.ResponseWriter.Header().Set(header, value)
	}
	log.Printf("[%d]: %s", status, w.Path)
	w.ResponseWriter.WriteHeader(status)
}
