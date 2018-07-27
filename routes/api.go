package routes

import (
	"github.com/gin-gonic/gin"
)

// DefineAPI
func DefineAPI(router gin.IRouter) {
	helloworldRouter := router.Group("/helloworld")
	DefineHelloworldRoute(helloworldRouter)

	userRouter := router.Group("/users")
	DefineUserRoute(userRouter)

	authRouter := router.Group("/auth")
	DefineAuthRoute(authRouter)

	blogpostRouter := router.Group("/blogposts")
	DefineBlogpostRoute(blogpostRouter)
}
