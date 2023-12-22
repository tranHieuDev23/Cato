package middlewares

import "net/http"

type inMemoryResponseWriter struct {
	Headers    http.Header
	StatusCode int
	Bodies     []byte
}

func newInMemoryResponseWriter() *inMemoryResponseWriter {
	return &inMemoryResponseWriter{
		Headers:    http.Header{},
		StatusCode: http.StatusOK,
		Bodies:     make([]byte, 0),
	}
}

func (w *inMemoryResponseWriter) Header() http.Header {
	return w.Headers
}

func (w *inMemoryResponseWriter) Write(data []byte) (int, error) {
	w.Bodies = append(w.Bodies, data...)
	return len(data), nil
}

func (w *inMemoryResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func (w *inMemoryResponseWriter) Apply(rw http.ResponseWriter) error {
	for header, value := range w.Headers {
		rw.Header()[header] = value
	}

	rw.WriteHeader(w.StatusCode)
	if _, err := rw.Write(w.Bodies); err != nil {
		return err
	}

	return nil
}
