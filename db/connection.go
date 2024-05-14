package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context, host, port, user, password, dbname string, isUnixSocket bool) (*pgx.Conn, error) {
	var connStr string
	if isUnixSocket {
		connStr = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	} else {
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	}
	return pgx.Connect(ctx, connStr)
}
