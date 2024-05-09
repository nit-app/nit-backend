package otp

import (
	"crypto/rand"
	"fmt"
	"github.com/nit-app/nit-backend/env"
	"sync"
)

type Generator interface {
	Generate() string
}

type secureRandomGenerator struct{}

func (srg *secureRandomGenerator) Generate() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%d%d%d%d", b[0]%10, b[1]%10, b[2]%10, b[3]%10)
}

type mockGenerator struct{}

func (mg *mockGenerator) Generate() string {
	return env.E().OtpMockAcceptCode
}

var (
	createOnce sync.Once
	gen        Generator
)

func createGenerator() {
	if len(env.E().OtpMockAcceptCode) != 0 {
		gen = &mockGenerator{}
	} else {
		gen = &secureRandomGenerator{}
	}
}

func Generate() string {
	createOnce.Do(createGenerator)
	return gen.Generate()
}
