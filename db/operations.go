package db

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FetchColumns function fetches column names for the given table from the source database.
func FetchColumns(ctx context.Context, pool *pgxpool.Pool, tableName string) ([]string, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT column_name FROM information_schema.columns WHERE table_name=$1", tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	return columns, nil
}

// DataProducer fetches data from the specified table and sends it to a channel.
func DataProducer(ctx context.Context, pool *pgxpool.Pool, tableName string, dataChan chan<- []interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(dataChan)

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Failed to acquire connection from pool: %v", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT * FROM "+tableName)
	if err != nil {
		log.Fatalf("Failed to fetch data from source database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatalf("Failed to read row values: %v", err)
		}
		dataChan <- values
	}
}

// DataConsumer receives data from a channel and writes it to the target database.
func DataConsumer(ctx context.Context, pool *pgxpool.Pool, tableName string, columns []string, dataChan <-chan []interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Failed to acquire connection from pool: %v", err)
	}
	defer conn.Release()

	rows := make([][]interface{}, 0)
	for row := range dataChan {
		rows = append(rows, row)
	}

	copyCount, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{tableName},
		columns,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Fatalf("Failed to copy data to target database: %v", err)
	}
	log.Printf("Copied %d rows to target database from table %s.", copyCount, tableName)
}
