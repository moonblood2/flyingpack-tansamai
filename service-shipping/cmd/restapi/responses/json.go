package responses

import (
	"encoding/json"
	"net/http"
)

// jsonResponse
// if ResponseWriter have Write() method jsonResponse have to implement it,
// by call Write() method of http.ResponseWriter again in Write() method of ResponseWriter.
type jsonResponseWriter struct {
	responseWriter http.ResponseWriter
}

// NewJSONResponseWriter create new ResponseWriter.
func NewJSONResponseWriter(w http.ResponseWriter) *jsonResponseWriter {
	return &jsonResponseWriter{
		responseWriter: w,
	}
}

// WriteWithStatus method for implement ResponseWriter interface, and write data in JSON with status code.
func (j jsonResponseWriter) WriteWithStatus(statusCode int, v interface{}) {
	j.responseWriter.Header().Set("Content-Type", "application/json")
	j.responseWriter.WriteHeader(statusCode)

	b, _ := json.Marshal(v)
	j.responseWriter.Write(b)
}
