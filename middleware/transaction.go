package middleware

import (
	"github.com/sinnott74/goblogserver/database"

	"github.com/gin-gonic/gin"
)

func Transaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := database.NewTransaction(c.Request.Context())
		ctx := database.SetTransaction(c.Request.Context(), t)
		c.Request = c.Request.WithContext(ctx)
		defer func() {
			if r := recover(); r != nil {
				err := t.Rollback()
				panic(err)
			} else if c.IsAborted() {
				err := t.Rollback()
				panic(err)
			} else {
				err := t.Commit()
				panic(err)
			}
		}()
		c.Next()
	}
}
