package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// StatusRecorder record status code and contain http.ResponseWriter
type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
	Byte       []byte
}

// WriteHeader call ResponseWriter WriteHeader() and save status code.
// Override http.ResponseWriter method for implementation.
func (r *StatusRecorder) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.StatusCode = statusCode
}

func (r *StatusRecorder) Write(b []byte) (int, error) {
	r.Byte = b
	return r.ResponseWriter.Write(b)
}

// Logging print log to standard output.
// If patch matched then mux will call handlerFunction.
// By the way, must wrapped every handles with logging.
// To solve this, create new mux for handle every path request with "/".
// Then call ServeHTTP of RootMux for dispatch to its multiplexer.
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[log]: Start: %v: [%v] %v\n", r.Method, r.RemoteAddr, r.RequestURI)

		//Request
		reqHeaderByte, _ := json.Marshal(r.Header)
		reqBodyByte, _ := ioutil.ReadAll(r.Body) //Copy body
		r.Body = ioutil.NopCloser(bytes.NewReader(reqBodyByte))
		log.Printf("Request : %v\n%v\n", string(reqHeaderByte), string(reqBodyByte))

		//Next
		started := time.Now()
		sr := &StatusRecorder{ResponseWriter: w}
		next.ServeHTTP(sr, r)

		//Response
		log.Printf("Response: %v\n", strings.TrimRight(string(sr.Byte), "\n"))

		log.Printf("[log]: End  : %v: [%v] %v [%v] [%vms]\n", r.RemoteAddr, r.Method, r.RequestURI, sr.StatusCode, time.Since(started).Milliseconds())
	})
}
