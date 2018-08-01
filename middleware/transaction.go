package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sinnott74/goblogserver/database"
)

func Transaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		t, ctx := database.NewTransaction(c.Request.Context())
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
