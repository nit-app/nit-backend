package env

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"sync"
)

const (
	OptionCarrierMethod  = "CARRIER_METHOD"
	OptionDatabaseString = "DATABASE_STRING"
	OptionRedisAddr      = "REDIS_ADDR"
	OptionRedisDB        = "REDIS_DB"
	OptionOtpAcceptCode  = "OTP_MOCK_ACCEPT_CODE"
)

var env *Env

var once sync.Once

func Init() {
	once.Do(func() {
		db := loadDB()
		redis := loadRedis()

		env = &Env{}
		env.DB = db
		env.Redis = redis

		env.CarrierMethod = os.Getenv(OptionCarrierMethod)
		env.OtpMockAcceptCode = os.Getenv(OptionOtpAcceptCode)
	})
	context.Background()
}
