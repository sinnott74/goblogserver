package main

import (
	"net/http"

	"github.com/sinnott74/goblogserver/database"
	"github.com/sinnott74/goblogserver/env"
	"github.com/sinnott74/goblogserver/routes"
)

func main() {

	db, err := database.Init()
	defer db.Close()
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":"+env.Port(), routes.Handler())
}
