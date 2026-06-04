package vodka

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
    http.ResponseWriter
    status      int
    wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

// Write intercepts implicit 200 responses. When a handler calls Write()
// directly — without a preceding explicit WriteHeader() call — the
// underlying http.ResponseWriter will implicitly send a 200 header.
// Without this method, rw.wroteHeader would stay false, causing
// Recovery() to incorrectly attempt a second WriteHeader(500) on an
// already-started response (corrupting the HTTP stream).
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.wroteHeader = true
		// status stays at its default 200 — the implicit write IS the 200.
	}
	return rw.ResponseWriter.Write(b)
}

func Logger() HandlerFunc {
	return func(c *Context) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: c.Writer, status: http.StatusOK}
		c.Writer = rw
		c.Next()
		log.Printf(
			Blue+"%s %s %d %s"+Reset,
			c.Request.Method,
			c.Request.URL.Path,
			rw.status,
			time.Since(start),
		)
	}
}
