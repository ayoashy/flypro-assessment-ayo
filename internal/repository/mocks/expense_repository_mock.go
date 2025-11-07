package mocks

import (
	"flypro-assessment-ayo/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockExpenseRepository struct {
	mock.Mock
}

func (m *MockExpenseRepository) Create(expense *models.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) GetByID(id uint) (*models.Expense, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expense), args.Error(1)
}

func (m *MockExpenseRepository) GetByUserID(userID uint, offset, limit int, category, status string) ([]models.Expense, int64, error) {
	args := m.Called(userID, offset, limit, category, status)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Expense), args.Get(1).(int64), args.Error(2)
}

func (m *MockExpenseRepository) Update(expense *models.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockExpenseRepository) GetByIDs(ids []uint) ([]models.Expense, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Expense), args.Error(1)
}

func (m *MockExpenseRepository) GetByUserIDAndIDs(userID uint, ids []uint) ([]models.Expense, error) {
	args := m.Called(userID, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Expense), args.Error(1)
}
