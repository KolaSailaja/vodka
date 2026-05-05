package vodka

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// httprouter wrapper
type Engine struct {
	router      *httprouter.Router
	middlewares []HandlerFunc
}

// creates a new router
func New() *Engine {
	return &Engine{
		router:      httprouter.New(),
		middlewares: make([]HandlerFunc, 0),
	}
}

func (e *Engine) Use(middleware ...HandlerFunc) {
	e.middlewares = append(e.middlewares, middleware...)
}

// Runs the http server
func (e *Engine) Run(addr string) error {
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Pouring Vodka on %s\n", addr)

	// Using net/http
	return http.ListenAndServe(addr, e.router)
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.router.GET(path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

		handlers := make([]HandlerFunc, 0, len(e.middlewares)+1)
		handlers = append(handlers, e.middlewares...)
		handlers = append(handlers, handler)

		c := &Context{
			Writer:   w,
			Request:  r,
			handlers: handlers,
			index:    -1,
		}

		c.Next()
	})
}
