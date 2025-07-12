package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kuzmindeniss/itk/internal/db/repository"
	"github.com/kuzmindeniss/itk/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) GetWalletByID(ctx context.Context, id uuid.UUID) (repository.Wallet, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repository.Wallet), args.Error(1)
}

func (m *MockWalletService) TopUpWalletBalance(ctx context.Context, id uuid.UUID, amount int32) (repository.Wallet, error) {
	args := m.Called(ctx, id, amount)
	return args.Get(0).(repository.Wallet), args.Error(1)
}

func setupTestRouter(mockService *MockWalletService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	handler := NewWalletHandler(mockService)

	r := gin.New()
	v1 := r.Group("/api/v1")
	v1.GET("/wallets/:id", handler.GetWallet)
	v1.POST("/wallet", handler.UpdateWalletBalance)

	return r
}

func TestWalletHandler_GetWallet_Success(t *testing.T) {
	mockService := new(MockWalletService)
	router := setupTestRouter(mockService)

	walletID := uuid.New()
	expectedWallet := repository.Wallet{
		ID:      walletID,
		Balance: 1000,
	}

	mockService.On("GetWalletByID", mock.Anything, walletID).Return(expectedWallet, nil)

	req, _ := http.NewRequest("GET", "/api/v1/wallets/"+walletID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response repository.Wallet
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedWallet.ID, response.ID)
	assert.Equal(t, expectedWallet.Balance, response.Balance)

	mockService.AssertExpectations(t)
}

func TestWalletHandler_GetWallet_InvalidID(t *testing.T) {
	mockService := new(MockWalletService)
	router := setupTestRouter(mockService)

	req, _ := http.NewRequest("GET", "/api/v1/wallets/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid wallet ID", response["error"])
}

func TestWalletHandler_GetWallet_ServiceError(t *testing.T) {
	mockService := new(MockWalletService)
	router := setupTestRouter(mockService)

	walletID := uuid.New()
	mockService.On("GetWalletByID", mock.Anything, walletID).Return(repository.Wallet{}, errors.New("database error"))

	req, _ := http.NewRequest("GET", "/api/v1/wallets/"+walletID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "database error", response["error"])

	mockService.AssertExpectations(t)
}

func TestWalletHandler_UpdateWalletBalance_Deposit_Success(t *testing.T) {
	mockService := new(MockWalletService)
	router := setupTestRouter(mockService)

	walletID := uuid.New()
	requestBody := UpdateBalanceRequest{
		Amount:        500,
		WalletID:      walletID.String(),
		OperationType: models.OperationDeposit,
	}

	expectedWallet := repository.Wallet{
		ID:      walletID,
		Balance: 1500,
	}

	mockService.On("TopUpWalletBalance", mock.Anything, walletID, int32(500)).Return(expectedWallet, nil)

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	wallet := response["wallet"].(map[string]interface{})
	assert.Equal(t, walletID.String(), wallet["id"])
	assert.Equal(t, float64(1500), wallet["balance"])

	mockService.AssertExpectations(t)
}

func TestWalletHandler_UpdateWalletBalance_Withdraw_Success(t *testing.T) {
	mockService := new(MockWalletService)
	router := setupTestRouter(mockService)

	walletID := uuid.New()
	requestBody := UpdateBalanceRequest{
		Amount:        300,
		WalletID:      walletID.String(),
		OperationType: models.OperationWithdraw,
	}

	expectedWallet := repository.Wallet{
		ID:      walletID,
		Balance: 700,
	}

	mockService.On("TopUpWalletBalance", mock.Anything, walletID, int32(-300)).Return(expectedWallet, nil)

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	wallet := response["wallet"].(map[string]interface{})
	assert.Equal(t, walletID.String(), wallet["id"])
	assert.Equal(t, float64(700), wallet["balance"])

	mockService.AssertExpectations(t)
}

func TestWalletHandler_UpdateWalletBalance_InvalidJSON(t *testing.T) {
	mockService := new(MockWalletService)
	router := setupTestRouter(mockService)

	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWalletHandler_UpdateWalletBalance_InvalidOperationType(t *testing.T) {
	mockService := new(MockWalletService)
	router := setupTestRouter(mockService)

	walletID := uuid.New()
	requestBody := UpdateBalanceRequest{
		Amount:        500,
		WalletID:      walletID.String(),
		OperationType: "INVALID",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid operation type", response["error"])
}

func TestWalletHandler_UpdateWalletBalance_InvalidWalletID(t *testing.T) {
	mockService := new(MockWalletService)
	router := setupTestRouter(mockService)

	requestBody := UpdateBalanceRequest{
		Amount:        500,
		WalletID:      "invalid-uuid",
		OperationType: models.OperationDeposit,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid wallet ID", response["error"])
}

func TestWalletHandler_UpdateWalletBalance_ServiceError(t *testing.T) {
	mockService := new(MockWalletService)
	router := setupTestRouter(mockService)

	walletID := uuid.New()
	requestBody := UpdateBalanceRequest{
		Amount:        500,
		WalletID:      walletID.String(),
		OperationType: models.OperationDeposit,
	}

	mockService.On("TopUpWalletBalance", mock.Anything, walletID, int32(500)).Return(repository.Wallet{}, errors.New("database error"))

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to update wallet balance", response["error"])

	mockService.AssertExpectations(t)
}
