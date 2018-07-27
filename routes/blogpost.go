package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sinnott74/goblogserver/model"
	"github.com/sinnott74/goblogserver/orm"
)

func DefineBlogpostRoute(router gin.IRouter) {
	router.GET("/", getBlogPosts)
}

func getBlogPosts(c *gin.Context) {
	ctx := c.Request.Context()
	blogposts := []model.BlogPost{}
	err := orm.SelectAll(ctx, &blogposts, &model.BlogPost{})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, blogposts)
}
