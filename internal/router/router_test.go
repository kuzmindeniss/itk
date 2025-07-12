package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kuzmindeniss/itk/internal/db/repository"
	"github.com/kuzmindeniss/itk/internal/handler"
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

func TestSetupRouter_RoutesRegistered(t *testing.T) {
	mockService := new(MockWalletService)
	walletHandler := handler.NewWalletHandler(mockService)

	router := SetupRouter(walletHandler)

	testCases := []struct {
		method   string
		path     string
		expected int
	}{
		{"GET", "/api/v1/wallets/invalid-uuid", http.StatusBadRequest},
		{"POST", "/api/v1/wallet", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest(tc.method, tc.path, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusNotFound, w.Code,
			"Route %s %s should be registered", tc.method, tc.path)
		assert.Equal(t, tc.expected, w.Code,
			"Route %s %s should return %d", tc.method, tc.path, tc.expected)
	}
}

func TestSetupRouter_CorrectRoutes(t *testing.T) {
	mockService := new(MockWalletService)
	walletHandler := handler.NewWalletHandler(mockService)
	router := SetupRouter(walletHandler)

	req, _ := http.NewRequest("GET", "/api/v1/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSetupRouter_APIVersion(t *testing.T) {
	mockService := new(MockWalletService)
	walletHandler := handler.NewWalletHandler(mockService)
	router := SetupRouter(walletHandler)

	testCases := []struct {
		path     string
		expected int
	}{
		{"/wallets/123", http.StatusNotFound},
		{"/wallet", http.StatusNotFound},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("GET", tc.path, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, tc.expected, w.Code,
			"Route %s should return %d", tc.path, tc.expected)
	}
}
