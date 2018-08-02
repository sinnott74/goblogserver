package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/gzip"
)

// Compression middleware gzip the response before senting to the client
func Compression() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}
