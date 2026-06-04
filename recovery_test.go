package vodka

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecovery_NoPanicWritesNormally(t *testing.T) {
	app := DefaultRouter()

	app.GET("/ok", func(c *Context) {
		c.JSON(http.StatusOK, M{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRecovery_PanicWrite(t *testing.T) {
	app := DefaultRouter()

	app.GET("/panic", func(c *Context) {
		panic("panic")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, Got %d", w.Code)
	}

}

func TestRecovery_PanicAfterWrite(t *testing.T) {
	app := DefaultRouter()

	app.GET("/panic-after-write", func(c *Context) {
		c.JSON(http.StatusOK, M{"data": "partial"})
		panic("panic after write")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic-after-write", nil)
	app.ServeHTTP(w, req)

	// Headers already sent, Recovery should NOT append a second body
	body := w.Body.String()
	if strings.Count(body, "{") > 1 {
		t.Errorf("expected single JSON body, got multiple concatenated: %s", body)
	}
}

// TestRecovery_PanicAfterRawWrite covers the case where a handler writes
// bytes directly via c.Writer.Write() — no explicit WriteHeader() call —
// and then panics. Before the Write() fix on responseWriter, rw.wroteHeader
// stayed false even though bytes were already in flight, so Recovery()
// would call c.JSON(500) a second time, corrupting the response stream.
func TestRecovery_PanicAfterRawWrite(t *testing.T) {
	app := DefaultRouter()

	app.GET("/raw-write-then-panic", func(c *Context) {
		// Bypass vodka helpers — write bytes directly.
		// The underlying ResponseWriter implicitly sends a 200 header here,
		// so rw.wroteHeader must be set to true by the Write() method.
		c.Writer.Write([]byte(`{"data":"partial"}`))
		panic("panic after raw write")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/raw-write-then-panic", nil)
	app.ServeHTTP(w, req)

	// Recovery must NOT append a 500 body on top of an already-started response.
	body := w.Body.String()
	if strings.Count(body, "{") > 1 {
		t.Errorf("expected single JSON body (no double-write), got: %s", body)
	}
}
