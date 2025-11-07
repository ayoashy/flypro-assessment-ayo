package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"flypro-assessment-ayo/internal/dto"
	"flypro-assessment-ayo/internal/models"
	repositorymocks "flypro-assessment-ayo/internal/repository/mocks"
	"flypro-assessment-ayo/internal/services"
	servicemocks "flypro-assessment-ayo/internal/services/mocks"
	"flypro-assessment-ayo/internal/utils"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExpenseService_CreateExpense(t *testing.T) {
	mockRepo := new(repositorymocks.MockExpenseRepository)
	mockCurrencyService := new(servicemocks.MockCurrencyService)
	mockRedis := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	service := services.NewExpenseService(mockRepo, mockCurrencyService, mockRedis, nil)

	tests := []struct {
		name          string
		input         dto.CreateExpenseRequest
		userID        uint
		setupMocks    func()
		expectedError error
		expectedID    uint
	}{
		{
			name: "valid expense creation",
			input: dto.CreateExpenseRequest{
				Amount:   100.50,
				Currency: "USD",
				Category: "travel",
				Description: "Flight ticket",
			},
			userID: 1,
			setupMocks: func() {
				mockRepo.On("Create", mock.AnythingOfType("*models.Expense")).Return(nil).Once()
				mockCurrencyService.On("ConvertCurrency", mock.Anything, 100.50, "USD", "USD").Return(100.50, nil).Maybe()
			},
			expectedID: 1,
		},
		{
			name: "invalid amount (zero)",
			input: dto.CreateExpenseRequest{
				Amount:   0,
				Currency: "USD",
				Category: "travel",
			},
			userID: 1,
			setupMocks: func() {
				// No mocks needed as validation should fail first
			},
			expectedError: errors.New("amount must be greater than 0"),
		},
		{
			name: "repository error",
			input: dto.CreateExpenseRequest{
				Amount:   100.50,
				Currency: "USD",
				Category: "travel",
			},
			userID: 1,
			setupMocks: func() {
				mockRepo.On("Create", mock.AnythingOfType("*models.Expense")).Return(errors.New("database error")).Once()
			},
			expectedError: errors.New("failed to create expense"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			mockCurrencyService.ExpectedCalls = nil
			mockCurrencyService.Calls = nil

			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			ctx := context.Background()
			result, err := service.CreateExpense(ctx, tt.userID, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedID > 0 {
					assert.Equal(t, tt.expectedID, result.ID)
				}
			}

			mockRepo.AssertExpectations(t)
			mockCurrencyService.AssertExpectations(t)
		})
	}
}

func TestExpenseService_GetExpenseByID(t *testing.T) {
	mockRepo := new(repositorymocks.MockExpenseRepository)
	mockCurrencyService := new(servicemocks.MockCurrencyService)
	mockRedis := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	service := services.NewExpenseService(mockRepo, mockCurrencyService, mockRedis, nil)

	tests := []struct {
		name          string
		expenseID     uint
		setupMocks    func()
		expectedError error
	}{
		{
			name:      "successful retrieval",
			expenseID: 1,
			setupMocks: func() {
				expense := &models.Expense{
					ID:          1,
					UserID:      1,
					Amount:      100.50,
					Currency:    "USD",
					Category:    "travel",
					Description: "Flight ticket",
					Status:      models.ExpenseStatusPending,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				mockRepo.On("GetByID", uint(1)).Return(expense, nil).Once()
			},
		},
		{
			name:      "expense not found",
			expenseID: 999,
			setupMocks: func() {
				mockRepo.On("GetByID", uint(999)).Return(nil, errors.New("record not found")).Once()
			},
			expectedError: utils.NewNotFoundError("expense"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			ctx := context.Background()
			result, err := service.GetExpenseByID(ctx, tt.expenseID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expenseID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestExpenseService_ListExpenses(t *testing.T) {
	mockRepo := new(repositorymocks.MockExpenseRepository)
	mockCurrencyService := new(servicemocks.MockCurrencyService)
	mockRedis := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	service := services.NewExpenseService(mockRepo, mockCurrencyService, mockRedis, nil)

	tests := []struct {
		name          string
		userID        uint
		filter        dto.ExpenseFilter
		setupMocks    func()
		expectedError error
		expectedCount int
	}{
		{
			name:   "successful list",
			userID: 1,
			filter: dto.ExpenseFilter{
				Page:    1,
				PerPage: 10,
			},
			setupMocks: func() {
				expenses := []models.Expense{
					{
						ID:          1,
						UserID:      1,
						Amount:      100.50,
						Currency:    "USD",
						Category:    "travel",
						Status:      models.ExpenseStatusPending,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				mockRepo.On("GetByUserID", uint(1), 0, 10, "", "").Return(expenses, int64(1), nil).Once()
			},
			expectedCount: 1,
		},
		{
			name:   "list with filter",
			userID: 1,
			filter: dto.ExpenseFilter{
				Page:     1,
				PerPage:  10,
				Category: "travel",
				Status:   "pending",
			},
			setupMocks: func() {
				expenses := []models.Expense{
					{
						ID:          1,
						UserID:      1,
						Amount:      100.50,
						Currency:    "USD",
						Category:    "travel",
						Status:      models.ExpenseStatusPending,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				mockRepo.On("GetByUserID", uint(1), 0, 10, "travel", "pending").Return(expenses, int64(1), nil).Once()
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			ctx := context.Background()
			result, err := service.ListExpenses(ctx, tt.userID, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedCount, len(result.Expenses))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}