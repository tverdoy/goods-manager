package transactor

import (
	"context"
	"database/sql"
)

type txKey struct{}

// Transactor represents a type that provides transaction management for database operations.
type Transactor struct {
	db *sql.DB
}

// Connection returns the current database connection.
// It takes a context.Context as a parameter.
// It returns a pointer to the current database transaction and a pointer to the underlying database connection.
func (t *Transactor) Connection(ctx context.Context) (*sql.Tx, *sql.DB) {
	return extractTx(ctx), t.db
}

// WithTransaction executes the provided function within a database transaction.
//
// It takes a context.Context and a function (fn) as parameters.
// The function fn is expected to perform database operations within the transaction.
// If fn returns an error, the transaction is rolled back, and the error is returned.
// If fn completes successfully, the transaction is committed.
//
// Example:
//
//	err := transactor.WithTransaction(ctx, func(txContext context.Context) error {
//	   return g.goodRepo.Create(ctx, good)
//	})
//	if err != nil {
//	    // Handle the error
//	}
func (t *Transactor) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	if err := fn(t.injectTx(ctx, tx)); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return tx.Commit()
}

// injectTx create new transaction and injects it into context
func (t *Transactor) injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}

	return nil
}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{db: db}
}
