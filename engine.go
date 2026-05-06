package vodka

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/julienschmidt/httprouter"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	engine      *Engine
}

// httprouter wrapper
type Engine struct {
	router *httprouter.Router
	*RouterGroup
}

// creates a new router
func NewRouter() *Engine {
	engine := &Engine{
		router: httprouter.New(),
	}

	engine.RouterGroup = &RouterGroup{
		prefix:      "",
		middlewares: make([]HandlerFunc, 0),
		engine:      engine,
	}

	return engine
}

func DefaultRouter() *Engine {
	engine := NewRouter()
	engine.Use(Recovery(), Logger())
	return engine
}

func (rg *RouterGroup) Group(prefix string, middlewares ...HandlerFunc) *RouterGroup {
	newMiddlewares := make([]HandlerFunc, len(rg.middlewares), len(rg.middlewares)+len(middlewares))
	copy(newMiddlewares, rg.middlewares)
	newMiddlewares = append(newMiddlewares, middlewares...)

	return &RouterGroup{
		prefix:      rg.prefix + prefix,
		middlewares: newMiddlewares,
		engine:      rg.engine,
	}
}

func (rg *RouterGroup) Use(middlewares ...HandlerFunc) {
	rg.middlewares = append(rg.middlewares, middlewares...)
}

// Runs the http server
func (e *Engine) Run(addr string) error {
	if addr == "" {
		addr = ":8080"
	}

	log.Printf(Green+"Pouring Vodka on %s\n"+Reset, addr)

	// Using net/http
	return http.ListenAndServe(addr, e.router)
}

// Serve Static files
func (rg *RouterGroup) Static(relativePath string, root string) {
	urlPattern := path.Join(relativePath, "/*filepath")

	fileServer := http.FileServer(http.Dir(root))

	rg.engine.router.GET(urlPattern, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		filepath := params.ByName("filepath")
		fullPath := path.Join(root, filepath)

		// Check if the requested file actually exists on the disk
		info, err := os.Stat(fullPath)

		// If the file doesn't exist OR it's a directory, serve index.html (React's entry point)
		if os.IsNotExist(err) || info.IsDir() {
			http.ServeFile(w, r, path.Join(root, "index.html"))
			return
		}

		// Otherwise, serve the actual file (css, js, images)
		// We use StripPrefix so /static/js/main.js looks in ./public/js/main.js
		http.StripPrefix(rg.prefix+relativePath, fileServer).ServeHTTP(w, r)
	})
}

func (rg *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	absolutePath := rg.prefix + comp

	handlers := make([]HandlerFunc, 0, len(rg.middlewares)+1)
	handlers = append(handlers, rg.middlewares...)
	handlers = append(handlers, handler)

	rg.engine.router.Handle(method, absolutePath, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		c := &Context{
			Writer:   w,
			Request:  r,
			Params:   params,
			handlers: handlers,
			index:    -1,
		}

		c.Next()
	})
}

func (rg *RouterGroup) GET(path string, handler HandlerFunc) {
	rg.addRoute(http.MethodGet, path, handler)
}

func (rg *RouterGroup) POST(path string, handler HandlerFunc) {
	rg.addRoute(http.MethodPost, path, handler)
}

func (rg *RouterGroup) PUT(path string, handler HandlerFunc) {
	rg.addRoute(http.MethodPut, path, handler)
}

func (rg *RouterGroup) DELETE(path string, handler HandlerFunc) {
	rg.addRoute(http.MethodDelete, path, handler)
}
