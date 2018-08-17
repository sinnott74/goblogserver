package orm

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type Transaction interface {
	// ORM
	Commit() error
	Rollback() error
	ID() string
	Tx() *sqlx.Tx
}

// type ORM interface {
// 	Execute(context.Context, string, ...interface{}) *sql.Row
// 	Save(context.Context, interface{}) error
// 	Insert(context.Context, interface{}) error
// 	Get(context.Context, interface{}) error
// 	Update(context.Context, interface{}, interface{}) error
// 	Delete(ctx context.Context, entity interface{}) error
// 	SelectAll(ctx context.Context, entities interface{}, where interface{}) error
// 	SelectOne(ctx context.Context, entity interface{}) error
// 	Count(ctx context.Context, where interface{}) (int64, error)
// }

// Config defines to options specfied at ORM initialisation
type Config struct {
	ConnURL          string
	DriverName       string
	MaxConns         int
	ToDBMapperFunc   func(string) string
	FromDBMapperFunc func(string) string
	Debug            bool
}

var (
	config         Config
	database       *sqlx.DB
	transactionKey = &contextKey{"Transaction"}
)

// Init initialised the database connection with configuration
func Init(c Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect(c.DriverName, c.ConnURL)
	if err != nil {
		return nil, err
	}
	if c.MaxConns != 0 {
		// db.SetMaxIdleConns(c.MaxConns)
		// db.SetMaxOpenConns(c.MaxConns)
	}
	db.MapperFunc(c.ToDBMapperFunc)
	config = c
	database = db
	return database, nil
}

// NewTransaction creates a new ORM transaction
func NewTransaction(ctx context.Context) (Transaction, context.Context) {
	if database == nil {
		panic(errors.New("ORM needs to be Initialised"))
	}
	t := &transactionImpl{uuid.NewV4(), database.MustBeginTx(ctx, nil)}
	c := setTransaction(ctx, t)
	return t, c
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "ORM context value " + k.name
}
