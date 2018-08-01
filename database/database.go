package database

import (
	"context"
	"fmt"
	"os"

	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // import postgres
	uuid "github.com/satori/go.uuid"
)

type key int

// TransactionKey is the key to access the Transaction object in a context
const TransactionKey key = 0

var db = initDB()

// Connect Initialise a DB connection
func initDB() *sqlx.DB {
	fmt.Println("Connecting to DB")
	connURL := connectionURL()
	database, err := sqlx.Open("postgres", connURL)
	if err != nil {
		panic(err)
	}
	err = database.Ping()
	if err != nil {
		panic(err)
	}
	database.SetMaxIdleConns(2)
	database.SetMaxOpenConns(2)

	// Map struct names to lower snake case
	database.MapperFunc(strcase.ToSnake)
	return database
}

// connectionURL get the database connection string from ENV Vars or used a default
func connectionURL() string {
	connectionString := os.Getenv("POSTGRES_URL")
	if connectionString == "" {
		connectionString = "postgres://Sinnott@localhost:5432/pwadb?sslmode=disable&timezone=UTC"
	}
	return connectionString
}

// Transaction controls access to the database
type Transaction interface {
	ID() uuid.UUID
	Commit() error
	Rollback() error
	Tx() *sqlx.Tx
}

// NewTransaction creates & begins a new database transaction
func NewTransaction(ctx context.Context) (Transaction, context.Context) {
	t := &transactionImpl{uuid.NewV4(), db.MustBeginTx(ctx, nil)}
	c := SetTransaction(ctx, t)
	return t, c
}

// Transaction struct
type transactionImpl struct {
	uuid uuid.UUID
	tx   *sqlx.Tx
}

// ID gets the transaction ID
func (t *transactionImpl) ID() uuid.UUID {
	return t.uuid
}

// Tx gets the database object
func (t *transactionImpl) Tx() *sqlx.Tx {
	return t.tx
}

// Commit a Transaction
func (t *transactionImpl) Commit() error {
	err := t.Tx().Commit()
	if err != nil {
		return err
	}
	fmt.Println("Commit " + t.ID().String())
	return nil
}

// Rollback a Transaction
func (t *transactionImpl) Rollback() error {
	err := t.Tx().Rollback()
	if err != nil {
		return err
	}
	fmt.Println("Rollback " + t.ID().String())
	return nil
}

// SetTransaction sets the Transaction in the context
func SetTransaction(ctx context.Context, t Transaction) context.Context {
	return context.WithValue(ctx, TransactionKey, t)
}

// GetTransaction gets the Transaction from the context
func GetTransaction(ctx context.Context) Transaction {
	return ctx.Value(TransactionKey).(Transaction)
}
