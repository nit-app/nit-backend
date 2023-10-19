package otp

import (
	"crypto/rand"
	"fmt"
	"github.com/nit-app/nit-backend/env"
)

type Generator interface {
	Generate() string
}

type secureRandomGenerator struct{}

func (srg *secureRandomGenerator) Generate() string {
	b := make([]byte, 3)

	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%d%d%d", b[0], b[1], b[2])
}

type mockGenerator struct{}

func (mg *mockGenerator) Generate() string {
	return env.E().OtpMockAcceptCode
}

func NewGenerator() Generator {
	if len(env.E().OtpMockAcceptCode) != 0 {
		return &mockGenerator{}
	}

	return &secureRandomGenerator{}
}
