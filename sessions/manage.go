package sessions

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/nit-app/nit-backend/env"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func getSession(token string) (*Session, error) {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tokHash := hash(token)

	fields, err := env.Redis().HGetAll(timeout, "session_"+tokHash).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if fields == nil {
		return nil, ErrNoSuchSession
	}

	s := &Session{}
	s.State = fields["state"]
	s.tokHash = tokHash

	if phone, ok := fields["phone"]; ok {
		attempt, _ := strconv.Atoi(fields["attempt"])

		s.OTP = &OtpState{
			Code:        fields["code"],
			Attempt:     attempt,
			PhoneNumber: phone,
		}
	}

	if s.State == StateAuthorized {
		sub := fields["sub"]
		s.Subject = &sub
		s.Authorized, _ = strconv.ParseInt(fields["iat"], 10, 64)
	}

	return s, nil
}

func createSession() (token string, err error) {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	token = genToken()
	err = env.Redis().HSet(timeout, "session_"+hash(token), "state", StateUnauthorized).Err()
	if err != nil {
		return
	}

	err = env.Redis().Expire(timeout, "session_"+hash(token), ttlSecs*time.Second).Err()

	return
}

func hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

const tokenPayloadSize = 64

func genToken() string {
	b := make([]byte, tokenPayloadSize)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}
