package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"flypro-assessment-ayo/internal/dto"
	"flypro-assessment-ayo/internal/handlers"
	"flypro-assessment-ayo/internal/services/mocks"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestExpenseHandler_CreateExpense(t *testing.T) {
	mockExpenseService := new(mocks.MockExpenseService)
	validator := validator.New()
	handler := handlers.NewExpenseHandler(mockExpenseService, validator)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func()
		expectedStatus int
	}{
		{
			name: "valid request",
			requestBody: dto.CreateExpenseRequest{
				Amount:   100.50,
				Currency: "USD",
				Category: "travel",
				Description: "Flight ticket",
			},
			setupMocks: func() {
				mockExpenseService.On("CreateExpense", mock.Anything, mock.AnythingOfType("uint"), mock.AnythingOfType("dto.CreateExpenseRequest")).
					Return(&dto.ExpenseResponse{ID: 1, Amount: 100.50, Currency: "USD"}, nil).Once()
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid request - missing amount",
			requestBody: map[string]interface{}{
				"currency": "USD",
				"category": "travel",
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid request - negative amount",
			requestBody: dto.CreateExpenseRequest{
				Amount:   -10,
				Currency: "USD",
				Category: "travel",
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExpenseService.ExpectedCalls = nil
			mockExpenseService.Calls = nil

			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			router := setupTestRouter()
			router.POST("/expenses", handler.CreateExpense)

			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/expenses", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", "1")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockExpenseService.AssertExpectations(t)
		})
	}
}