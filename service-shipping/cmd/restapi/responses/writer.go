package responses

// ResponseWriter interface
type ResponseWriter interface {
	WriteWithStatus(statusCode int, v interface{})
}
