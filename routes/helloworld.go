package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DefineHelloworldRoute(router gin.IRouter) {
	router.GET("/", getHelloWorld)
}

func getHelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello world")
}
