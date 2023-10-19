package sms

import "go.uber.org/zap"

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

func NewCarrier() Carrier {
	return &mockCarrier{}
}
