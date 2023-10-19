package env

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

func loadDB() *sql.DB {
	databaseString := os.Getenv(OptionDatabaseString)
	d, err := sql.Open("postgres", databaseString)

	if err != nil {
		panic(err)
	}

	return d
}

func loadRedis() *redis.Client {
	dbIndex, _ := strconv.Atoi(os.Getenv(OptionRedisDB))

	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv(OptionRedisAddr),
		DB:   dbIndex,
	})

	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(timeout).Err(); err != nil {
		panic(err)
	}

	return client
}
