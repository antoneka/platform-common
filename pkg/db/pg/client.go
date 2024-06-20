package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/antoneka/platform-common/pkg/db"
)

// pgClient represents the implementation of the database client interface.
type pgClient struct {
	masterDBC db.DB
}

// NewDBClient creates a new PostgreSQL client instance.
func NewDBClient(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}

	return &pgClient{
		masterDBC: NewDB(dbc),
	}, nil
}

// DB returns the master database connection.
func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

// Close closes the database connection.
func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
