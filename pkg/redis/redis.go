package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func Connect() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		log.Printf("Invalid REDIS_DB value, defaulting to 0. Error: %v", err)
		db = 0
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis!")
	return rdb
}