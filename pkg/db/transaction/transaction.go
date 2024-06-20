package transaction

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/antoneka/auth/pkg/client/db"
	"github.com/antoneka/auth/pkg/client/db/pg"
)

// manager implements the TxManager interface.
type manager struct {
	db db.Transactor
}

// NewTransactionManager creates a new instance of the transaction manager.
func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

// transaction performs the transactional logic.
func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, f db.Handler) (err error) {
	// Check if a transaction is already in progress in the context.
	// If a transaction is already in progress, skip initializing a new transaction
	// and execute the handler directly.
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return f(ctx)
	}

	// If there is no transaction in progress, begin a new one.
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin the transaction: %w", err)
	}

	// Update the context with the transaction.
	ctx = pg.MakeContextTx(ctx, tx)

	// Defer the rollback or commit operation based on the outcome of the transactional logic.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}

		if err != nil {
			// If an error occurred, rollback the transaction.
			if errRollBack := tx.Rollback(ctx); errRollBack != nil {
				err = fmt.Errorf("tx rollback failed: %v: %w", errRollBack, err)
			}
		} else {
			// If no error occurred, commit the transaction.
			if err = tx.Commit(ctx); err != nil {
				err = fmt.Errorf("tx commit failed: %w", err)
			}
		}
	}()

	// Execute the transactional logic.
	// If the handler returns an error, rollback the transaction, otherwise commit the transaction.
	if err = f(ctx); err != nil {
		err = fmt.Errorf("failed to execute the code inside the transaction: %w", err)
	}

	return err
}

// ReadCommitted executes the provided handler within a transaction with the Read Committed isolation level.
func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

	return m.transaction(ctx, txOpts, f)
}
