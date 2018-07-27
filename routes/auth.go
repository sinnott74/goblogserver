package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sinnott74/goblogserver/model"
	"github.com/sinnott74/goblogserver/orm"
)

func DefineAuthRoute(router gin.IRouter) {
	router.POST("/signup", signUp)
	router.POST("/login", login)
}

type signUpRequest struct {
	model.User
	model.Credential
}

func signUp(c *gin.Context) {
	ctx := c.Request.Context()
	var signUpRequest signUpRequest
	c.BindJSON(&signUpRequest)
	orm.Insert(ctx, &signUpRequest.User)
	signUpRequest.UserID = signUpRequest.User.ID
	orm.Insert(ctx, &signUpRequest.Credential)
	fmt.Printf("%+v\n", signUpRequest)
}

type loginRequest struct {
	Username string
	Password string
}

func login(c *gin.Context) {
	ctx := c.Request.Context()
	var loginRequest loginRequest
	c.BindJSON(&loginRequest)
	credential := &model.Credential{}
	authenticated := credential.Authenticate(ctx, loginRequest.Username, loginRequest.Password)
	c.JSON(http.StatusOK, authenticated)
}
