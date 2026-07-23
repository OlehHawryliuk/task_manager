package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectToRedis() *redis.Client {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := os.Getenv("REDIS_PASSWORD")
	poolSize, _ := strconv.Atoi(getEnv("REDIS_POOL_SIZE", "50"))

	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", host, port),
		Password:        password,
		DB:              0,
		PoolSize:        poolSize,
		MinIdleConns:    10,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		MaxRetries:      3,
		MinRetryBackoff: 50 * time.Millisecond,
	})

	var pong string
	var err error
	ctx := context.Background()
	for i := 1; i <= 3; i++ {
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		pong, err = client.Ping(pingCtx).Result()
		cancel()

		if err == nil {
			break
		}

		log.Printf("Attempt %d: Redis unavailable, waiting... (%v)\n", i, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Printf("Fatal: Could not connect to Redis: %v\n", err)
		return nil
	}

	log.Printf("Redis connected successfully: %s (pool: %d)\n", pong, poolSize)
	RedisClient = client
	return client
}
