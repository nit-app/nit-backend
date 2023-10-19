package env

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

func DB() *sql.DB {
	return env.DB
}

func Redis() *redis.Client {
	return env.Redis
}

func E() *Env {
	return env
}
