package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// GenerateConnString constructs a PostgreSQL connection string.
func GenerateConnString(host, port, user, password, dbname string, isUnixSocket bool) string {
	if isUnixSocket {
		return fmt.Sprintf("postgres://%s:%s@/%s?host=%s", user, password, dbname, host)
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
}

// Connect establishes a connection to the database using a connection string.
func Connect(ctx context.Context, host, port, user, password, dbname string, isUnixSocket bool) (*pgx.Conn, error) {
	connStr := GenerateConnString(host, port, user, password, dbname, isUnixSocket)
	return pgx.Connect(ctx, connStr)
}
