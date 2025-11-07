package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockCurrencyService struct {
	mock.Mock
}

func (m *MockCurrencyService) ConvertCurrency(ctx context.Context, amount float64, from, to string) (float64, error) {
	args := m.Called(ctx, amount, from, to)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockCurrencyService) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(float64), args.Error(1)
}
