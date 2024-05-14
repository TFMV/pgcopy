package db

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
)

// DataProducer fetches data from the specified table and sends it to a channel.
func DataProducer(ctx context.Context, conn *pgx.Conn, tableName string, dataChan chan<- []interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(dataChan)

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
func DataConsumer(ctx context.Context, conn *pgx.Conn, tableName string, columns []string, dataChan <-chan []interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	var data [][]interface{}
	for row := range dataChan {
		data = append(data, row)
	}

	copyCount, err := conn.CopyFrom(ctx, pgx.Identifier{tableName}, columns, pgx.CopyFromSlice(len(data), func(i int) ([]interface{}, error) {
		return data[i], nil
	}))
	if err != nil {
		log.Fatalf("Failed to copy data to target database: %v", err)
	}
	log.Printf("Copied %d rows to target database from table %s.", copyCount, tableName)
}
