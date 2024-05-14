package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/TFMV/pgcopy/db"
	"github.com/TFMV/pgcopy/model"
)

func main() {
	config, err := model.GetConf("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	ctx := context.Background()
	var wg sync.WaitGroup
	startTime := time.Now()

	sourcePool, err := db.ConnectPool(ctx, config.Source.Host, config.Source.Port, config.Source.User, config.Source.Pass, config.Source.DB, config.Source.IsUnixSocket)
	if err != nil {
		log.Fatalf("Unable to connect to source database: %v", err)
	}
	defer sourcePool.Close()

	targetPool, err := db.ConnectPool(ctx, config.Target.Host, config.Target.Port, config.Target.User, config.Target.Pass, config.Target.DB, config.Target.IsUnixSocket)
	if err != nil {
		log.Fatalf("Unable to connect to target database: %v", err)
	}
	defer targetPool.Close()

	for _, tableName := range config.Tables {
		columns, err := db.FetchColumns(ctx, sourcePool, tableName)
		if err != nil {
			log.Fatalf("Failed to fetch columns from source database for table %s: %v", tableName, err)
		}

		dataChan := make(chan []interface{}, 10000)

		wg.Add(2)
		go db.DataProducer(ctx, sourcePool, tableName, dataChan, &wg)
		go db.DataConsumer(ctx, targetPool, tableName, columns, dataChan, &wg)
	}

	wg.Wait()

	elapsedTime := time.Since(startTime)
	response := model.JsonResponse{
		Message:   "Data replication completed successfully",
		TimeTaken: elapsedTime.Seconds(),
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error encoding JSON response: %v", err)
	}

	fmt.Printf("Response: %s\n", string(respBytes))
	fmt.Printf("Time taken: %.2f seconds\n", elapsedTime.Seconds())

	fmt.Println(string(respBytes))
}
