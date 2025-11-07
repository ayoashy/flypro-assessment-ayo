package mocks

import (
	"context"

	"flypro-assessment-ayo/internal/dto"

	"github.com/stretchr/testify/mock"
)

type MockExpenseService struct {
	mock.Mock
}

func (m *MockExpenseService) CreateExpense(ctx context.Context, userID uint, req dto.CreateExpenseRequest) (*dto.ExpenseResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ExpenseResponse), args.Error(1)
}

func (m *MockExpenseService) GetExpenseByID(ctx context.Context, id uint) (*dto.ExpenseResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ExpenseResponse), args.Error(1)
}

func (m *MockExpenseService) ListExpenses(ctx context.Context, userID uint, filter dto.ExpenseFilter) (*dto.ExpenseListResponse, error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ExpenseListResponse), args.Error(1)
}

func (m *MockExpenseService) UpdateExpense(ctx context.Context, id, userID uint, req dto.UpdateExpenseRequest) (*dto.ExpenseResponse, error) {
	args := m.Called(ctx, id, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ExpenseResponse), args.Error(1)
}

func (m *MockExpenseService) DeleteExpense(ctx context.Context, id, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}