package database

import (
	"github.com/iancoleman/strcase"
	_ "github.com/lib/pq" // import postgres
	"github.com/sinnott74/goblogserver/env"
	"github.com/sinnott74/goblogserver/orm"
)

func Init() error {
	config := orm.Config{
		ConnURL:          env.ConnectionURL(),
		DriverName:       "postgres",
		MaxConns:         2,
		ToDBMapperFunc:   strcase.ToSnake,
		FromDBMapperFunc: strcase.ToCamel,
	}
	return orm.Init(config)
}
