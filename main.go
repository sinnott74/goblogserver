package main

import (
	"os"

	"github.com/sinnott74/goblogserver/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/sinnott74/goblogserver/database"
	"github.com/sinnott74/goblogserver/routes"
)

func main() {

	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.RedirectToHTTPSRouter())
	r.Use(gin.Logger())
	r.Use(middleware.Compression())
	r.Use(middleware.ETag())
	r.Use(gin.Recovery())
	r.Use(middleware.Transaction())

	apiRouter := r.Group("/api")
	routes.DefineAPI(apiRouter)

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// s := &http.Server{
	// 	Addr:    ":" + port,
	// 	Handler: r,
	// }
	// s.SetKeepAlivesEnabled

	// s.ListenAndServe()

	// Listen for requests
	r.Run(":" + port)
}
