package middleware

import (
	"github.com/gin-gonic/gin"
)

func ETag() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Connection", "keep-alive")
		// defer func() {
		// 	c.Header("Etag", "33a64df551425fcc55e4d42a148795d9f25f89d4")
		// }()
		c.Next()
	}
}
