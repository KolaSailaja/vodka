package vodka

import (
	"log"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		log.Printf("%s %s", c.Request.Method, c.Request.URL.Path)
	}
}
