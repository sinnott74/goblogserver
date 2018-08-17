package database

import (
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq" // import postgres
	"github.com/sinnott74/goblogserver/env"
	"github.com/sinnott74/goblogserver/orm"
)

// Init intialises a database connection
func Init() (*sqlx.DB, error) {
	config := orm.Config{
		ConnURL:          env.ConnectionURL(),
		DriverName:       "postgres",
		MaxConns:         2,
		ToDBMapperFunc:   strcase.ToSnake,
		FromDBMapperFunc: strcase.ToCamel,
		Debug:            env.Debug(),
	}
	return orm.Init(config)
}
