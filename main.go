package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/TFMV/pgcopy/db"
)

type Config struct {
	Source struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		User         string `yaml:"user"`
		Pass         string `yaml:"pass"`
		DB           string `yaml:"db"`
		IsUnixSocket bool   `yaml:"isUnixSocket"`
	} `yaml:"source"`

	Target struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		User         string `yaml:"user"`
		Pass         string `yaml:"pass"`
		DB           string `yaml:"db"`
		IsUnixSocket bool   `yaml:"isUnixSocket"`
	} `yaml:"target"`

	Tables []string `yaml:"tables"`
}

type JsonResponse struct {
	Message   string        `json:"message"`
	TimeTaken time.Duration `json:"timeTaken"`
}

func (c *Config) getConf() *Config {
	yamlContent, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlContent, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func main() {
	var config Config
	config.getConf()

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
	response := JsonResponse{
		Message:   "Data replication completed successfully",
		TimeTaken: elapsedTime,
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error encoding JSON response: %v", err)
	}

	fmt.Println(string(respBytes))
}
