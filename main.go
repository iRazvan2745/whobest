package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func createEnvFile() {
	envFileName := ".env"
	if _, err := os.Stat(envFileName); os.IsNotExist(err) {
		file, err := os.Create(envFileName)
		if err != nil {
			log.Fatalf("Error creating .env file: %v", err)
		}
		defer file.Close()

		defaultContent := `DRAGONFLY_ADDR=localhost:6379
KEYDB_ADDR=localhost:6380
REDIS_ADDR=localhost:6381
REDIS_PASSWORD=your_password_here
`
		_, err = file.WriteString(defaultContent)
		if err != nil {
			log.Fatalf("Error writing to .env file: %v", err)
		}
	}
}

func loadConfig() (string, string, string) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return os.Getenv("DRAGONFLY_ADDR"),
		os.Getenv("KEYDB_ADDR"),
		os.Getenv("REDIS_ADDR")
}

func testDB(client *redis.Client, dbName string, duration time.Duration) (int64, int64) {
	var ops, errors int64
	var wg sync.WaitGroup

	endTime := time.Now().Add(duration)
	concurrencyLevel := 1000

	wg.Add(concurrencyLevel)
	for i := 0; i < concurrencyLevel; i++ {
		go func(workerID int) {
			defer wg.Done()

			pipe := client.Pipeline()
			defer pipe.Close()

			keyBuf := make([]byte, 32)
			valueBuf := []byte("value")
			hashKey := make([]byte, 37)
			listKey := make([]byte, 37)

			for time.Now().Before(endTime) {
				n := copy(keyBuf, fmt.Sprintf("k%d%d", workerID, time.Now().UnixNano()))
				key := keyBuf[:n]

				copy(hashKey, key)
				copy(hashKey[n:], ":hash")
				copy(listKey, key)
				copy(listKey[n:], ":list")

				pipe.Set(ctx, string(key), valueBuf, 0)
				pipe.HSet(ctx, string(hashKey), "field1", "value1", "field2", "value2")
				pipe.LPush(ctx, string(listKey), "item1", "item2", "item3")
				pipe.Get(ctx, string(key))

				cmds, err := pipe.Exec(ctx)
				if err != nil {
					atomic.AddInt64(&errors, 1)
					continue
				}

				for _, cmd := range cmds {
					if cmd.Err() != nil {
						atomic.AddInt64(&errors, 1)
					} else {
						atomic.AddInt64(&ops, 1)
					}
				}

				pipe.Discard()
			}
		}(i)
	}

	wg.Wait()
	return ops, errors
}

func main() {
	results := make([]string, 0, 3)

	createEnvFile()

	dragonflyAddr, keydbAddr, redisAddr := loadConfig()

	baseOpts := &redis.Options{
		PoolSize:           10000,
		MinIdleConns:       2000,
		MaxRetries:         0,
		ReadTimeout:        25 * time.Millisecond,
		WriteTimeout:       25 * time.Millisecond,
		PoolTimeout:        15 * time.Second,
		IdleTimeout:        15 * time.Second,
		IdleCheckFrequency: 15 * time.Second,
	}

	dragonflyOpts := *baseOpts
	dragonflyOpts.Addr = dragonflyAddr
	dragonflyClient := redis.NewClient(&dragonflyOpts)
	defer func() {
		if err := dragonflyClient.Close(); err != nil {
			log.Fatalf("Error closing DragonflyDB client: %v", err)
		}
	}()

	keydbOpts := *baseOpts
	keydbOpts.Addr = keydbAddr
	keydbClient := redis.NewClient(&keydbOpts)
	defer func() {
		if err := keydbClient.Close(); err != nil {
			log.Fatalf("Error closing KeyDB client: %v", err)
		}
	}()

	redisOpts := *baseOpts
	redisOpts.Addr = redisAddr
	redisClient := redis.NewClient(&redisOpts)
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatalf("Error closing Redis client: %v", err)
		}
	}()

	log.Println("Testing DragonflyDB...")
	dragonflyOps, dragonflyErrors := testDB(dragonflyClient, "DragonflyDB", 5*time.Second)
	results = append(results, fmt.Sprintf(
		"DragonflyDB:\n"+
			"  - Total Operations: %d\n"+
			"  - Errors: %d\n"+
			"  - Successful Operations: %d\n",
		dragonflyOps, dragonflyErrors, dragonflyOps-dragonflyErrors))

	log.Println("Testing KeyDB...")
	keydbOps, keydbErrors := testDB(keydbClient, "KeyDB", 5*time.Second)
	results = append(results, fmt.Sprintf(
		"KeyDB:\n"+
			"  - Total Operations: %d\n"+
			"  - Errors: %d\n"+
			"  - Successful Operations: %d\n",
		keydbOps, keydbErrors, keydbOps-keydbErrors))

	log.Println("Testing Redis...")
	redisOps, redisErrors := testDB(redisClient, "Redis", 5*time.Second)
	results = append(results, fmt.Sprintf(
		"Redis:\n"+
			"  - Total Operations: %d\n"+
			"  - Errors: %d\n"+
			"  - Successful Operations: %d\n",
		redisOps, redisErrors, redisOps-redisErrors))

	file, err := os.Create("results.txt")
	if err != nil {
		log.Fatalf("Error creating results file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("Error closing results file: %v", err)
		}
	}()

	content := "Performance Test Results:\n\n"
	for _, result := range results {
		content += result + "\n"
	}

	if _, err := file.WriteString(content); err != nil {
		log.Fatalf("Error writing to results file: %v", err)
	}
}
