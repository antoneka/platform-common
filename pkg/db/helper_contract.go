package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Tx represents a database transaction.
// This interface is used for transaction testing.
type Tx interface {
	// Begin starts a pseudo nested transaction.
	Begin(ctx context.Context) (pgx.Tx, error)

	// BeginFunc starts a pseudo nested transaction and executes f. If f does not return an err the pseudo nested
	// transaction will be committed. If it does then it will be rolled back.
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) (err error)

	// Commit commits the transaction if this is a real transaction or releases the savepoint if this is a pseudo nested
	// transaction. Commit will return ErrTxClosed if the Tx is already closed, but is otherwise safe to call multiple
	// times. If the commit fails with a rollback status (e.g. the transaction was already in a broken state) then
	// ErrTxCommitRollback will be returned.
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction if this is a real transaction or rolls back to the savepoint if this is a
	// pseudo nested transaction. Rollback will return ErrTxClosed if the Tx is already closed, but is otherwise safe to
	// call multiple times. Hence, a defer tx.Rollback() is safe even if tx.Commit() will be called first in a non-error
	// condition. Any other failure of a real transaction will result in the connection being closed.
	Rollback(ctx context.Context) error

	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	LargeObjects() pgx.LargeObjects

	Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error)

	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)

	// Conn returns the underlying *Conn that on which this transaction is executing.
	Conn() *pgx.Conn
}
