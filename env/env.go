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
	ListenAddress     string
}

const defaultListenAddress = ":8080"

func (e *Env) Validate() {
	if e.DB == nil {
		panic("database connection must be open")
	}

	if len(e.ListenAddress) == 0 {
		e.ListenAddress = defaultListenAddress
	}
}
