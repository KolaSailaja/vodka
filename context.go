package vodka

import (
	"encoding/json"
	"net/http"
)

type HandlerFunc func(*Context) // Handler Function with Context wrapping

type Context struct {
	Writer   http.ResponseWriter // net/http response writer
	Request  *http.Request       // net/http request
	handlers []HandlerFunc       // stores middleware funcs and also main handler func
	index    int8                // tracks current step
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) JSON(statusCode int, obj any) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(statusCode)
	json.NewEncoder(c.Writer).Encode(obj)
}

func (c *Context) String(statusCode int, text string) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(statusCode)
	c.Writer.Write([]byte(text))
}
