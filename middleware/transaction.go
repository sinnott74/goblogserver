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
			r := recover()
			if r != nil {
				t.Rollback()
				panic(r)
			} else {
				t.Commit()
			}
		}()
		c.Next()
	}
}
