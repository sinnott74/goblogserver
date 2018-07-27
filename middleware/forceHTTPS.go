package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RedirectToHTTPSRouter is middleware which redirects the user to https if the x-forward-proto header is set to http
func RedirectToHTTPSRouter() gin.HandlerFunc {
	return func(c *gin.Context) {
		proto := c.Request.Header.Get("x-forwarded-proto")
		// fmt.Printf("%+v\n", req)
		if proto == "http" || proto == "HTTP" {
			c.Redirect(http.StatusPermanentRedirect, "https://"+c.Request.Host+c.Request.RequestURI)
			return
		}
		c.Next()
	}
}
