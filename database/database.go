package database

import (
	"os"

	"github.com/iancoleman/strcase"
	_ "github.com/lib/pq" // import postgres
	"github.com/sinnott74/goblogserver/orm"
)

func Init() error {
	config := orm.Config{
		ConnURL:          connectionURL(),
		DriverName:       "postgres",
		MaxConns:         2,
		ToDBMapperFunc:   strcase.ToSnake,
		FromDBMapperFunc: strcase.ToCamel,
	}
	return orm.Init(config)
}

// connectionURL get the database connection string from ENV Vars or used a default
func connectionURL() string {
	connectionString := os.Getenv("POSTGRES_URL")
	if connectionString == "" {
		connectionString = "postgres://Sinnott@localhost:5432/pwadb?sslmode=disable&timezone=UTC"
	}
	return connectionString
}
