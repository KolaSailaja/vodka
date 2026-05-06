package vodka

import "net/http"

func AllowCORS(origins []string) HandlerFunc {
	return func(c *Context) {
		clientOrigin := c.Request.Header.Get("Origin")
		isAllowed := false

		for _, o := range origins {
			if o == "*" || o == clientOrigin {
				isAllowed = true

				c.Writer.Header().Set("Access-Control-Allow-Origin", clientOrigin)
				c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

				break
			}
		}

		if c.Request.Method == http.MethodOptions {
			if isAllowed {
				c.Writer.WriteHeader(http.StatusNoContent) // 204: All good, proceed!
			} else {
				c.Writer.WriteHeader(http.StatusForbidden) // 403: Origin not allowed
			}

			c.Abort()
			return
		}

		c.Next()
	}
}
