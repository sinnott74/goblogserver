package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sinnott74/goblogserver/model"
	"github.com/sinnott74/goblogserver/orm"
)

func DefineUserRoute(router gin.IRouter) {
	router.GET("/", getUsers)
	router.GET("/:id", getUser)
	router.POST("/", createUser)
}

func getUsers(c *gin.Context) {
	ctx := c.Request.Context()
	users := []model.User{}
	orm.SelectAll(ctx, &users, &model.User{})
	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 0)
	if err != nil {
		panic(err)
	}
	ctx := c.Request.Context()
	user := model.User{ID: id}
	orm.Get(ctx, &user)
	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	var user model.User
	c.BindJSON(&user)
	ctx := c.Request.Context()
	orm.Insert(ctx, &user)
	c.JSON(http.StatusOK, user)
}
