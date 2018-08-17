package orm

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

// transactionImpl is the implentation of the Transaction interface
type transactionImpl struct {
	uuid uuid.UUID
	tx   *sqlx.Tx
}

// ID of a Transaction
func (t *transactionImpl) ID() string {
	return t.uuid.String()
}

// Commit a Transaction
func (t *transactionImpl) Commit() error {
	if config.Debug {
		fmt.Println("Commit " + t.ID())
	}
	return t.tx.Commit()
}

// Rollback a Transaction
func (t *transactionImpl) Rollback() error {
	if config.Debug {
		fmt.Println("Rollback " + t.ID())
	}
	return t.tx.Rollback()
}

func (t *transactionImpl) Tx() *sqlx.Tx {
	return t.tx
}

// setTransaction sets the Transaction on the context
func setTransaction(ctx context.Context, t Transaction) context.Context {
	return context.WithValue(ctx, transactionKey, t)
}

// getTransaction gets the Transaction from the context
func getTransaction(ctx context.Context) Transaction {
	return ctx.Value(transactionKey).(Transaction)
}
