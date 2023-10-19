package env

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type Env struct {
	DB                *sql.DB
	Redis             *redis.Client
	CarrierMethod     string
	OtpMockAcceptCode string
}

func (e *Env) Validate() {
	if e.DB == nil {
		panic("database connection must be open")
	}
}
