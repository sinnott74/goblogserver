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
	err := t.tx.Commit()
	if err != nil {
		return err
	}
	fmt.Println("Commit " + t.ID())
	return nil
}

// Rollback a Transaction
func (t *transactionImpl) Rollback() error {
	err := t.tx.Rollback()
	if err != nil {
		return err
	}
	fmt.Println("Rollback " + t.ID())
	return nil
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
