package sms

import (
	"go.uber.org/zap"
	"sync"
)

const (
	CarrierMethodMock      = "mock"
	CarrierMethodSmsprosto = "smsprosto"
)

var Methods = map[string]struct{}{CarrierMethodMock: {}, CarrierMethodSmsprosto: {}}

type Carrier interface {
	Send(phoneNumber string, text string) error
}

type mockCarrier struct{}

func (m *mockCarrier) Send(phoneNumber string, text string) error {
	zap.S().Infow("sending sms", "phoneNumber", phoneNumber, "text", text)

	return nil
}

var (
	createOnce sync.Once
	car        Carrier
)

func createCarrier() {
	car = &mockCarrier{}
}

func Send(phoneNumber string, text string) error {
	createOnce.Do(createCarrier)
	return car.Send(phoneNumber, text)
}
