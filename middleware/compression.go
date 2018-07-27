package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/gzip"
)

func Compression() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}
