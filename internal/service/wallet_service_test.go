package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/kuzmindeniss/itk/internal/db/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetWalletByID(ctx context.Context, id uuid.UUID) (repository.Wallet, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repository.Wallet), args.Error(1)
}

func (m *MockRepository) UpdateWallet(ctx context.Context, arg repository.UpdateWalletParams) (repository.Wallet, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(repository.Wallet), args.Error(1)
}

func TestWalletService_GetWalletByID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewWalletService(mockRepo)

	ctx := context.Background()
	walletID := uuid.New()
	expectedWallet := repository.Wallet{
		ID:      walletID,
		Balance: 1000,
	}

	mockRepo.On("GetWalletByID", ctx, walletID).Return(expectedWallet, nil)

	result, err := service.GetWalletByID(ctx, walletID)

	assert.NoError(t, err)
	assert.Equal(t, expectedWallet.ID, result.ID)
	assert.Equal(t, expectedWallet.Balance, result.Balance)

	mockRepo.AssertExpectations(t)
}

func TestWalletService_GetWalletByID_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewWalletService(mockRepo)

	ctx := context.Background()
	walletID := uuid.New()
	expectedError := errors.New("wallet not found")

	mockRepo.On("GetWalletByID", ctx, walletID).Return(repository.Wallet{}, expectedError)

	result, err := service.GetWalletByID(ctx, walletID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, repository.Wallet{}, result)

	mockRepo.AssertExpectations(t)
}

func TestWalletService_TopUpWalletBalance_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewWalletService(mockRepo)

	ctx := context.Background()
	walletID := uuid.New()
	amount := int32(500)

	expectedParams := repository.UpdateWalletParams{
		ID:     walletID,
		Amount: amount,
	}

	expectedWallet := repository.Wallet{
		ID:      walletID,
		Balance: 1500,
	}

	mockRepo.On("UpdateWallet", ctx, expectedParams).Return(expectedWallet, nil)

	result, err := service.TopUpWalletBalance(ctx, walletID, amount)

	assert.NoError(t, err)
	assert.Equal(t, expectedWallet.ID, result.ID)
	assert.Equal(t, expectedWallet.Balance, result.Balance)

	mockRepo.AssertExpectations(t)
}

func TestWalletService_TopUpWalletBalance_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewWalletService(mockRepo)

	ctx := context.Background()
	walletID := uuid.New()
	amount := int32(500)

	expectedParams := repository.UpdateWalletParams{
		ID:     walletID,
		Amount: amount,
	}

	expectedError := errors.New("database update failed")

	mockRepo.On("UpdateWallet", ctx, expectedParams).Return(repository.Wallet{}, expectedError)

	result, err := service.TopUpWalletBalance(ctx, walletID, amount)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, repository.Wallet{}, result)

	mockRepo.AssertExpectations(t)
}

func TestWalletService_TopUpWalletBalance_Withdraw(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewWalletService(mockRepo)

	ctx := context.Background()
	walletID := uuid.New()
	amount := int32(-300)

	expectedParams := repository.UpdateWalletParams{
		ID:     walletID,
		Amount: amount,
	}

	expectedWallet := repository.Wallet{
		ID:      walletID,
		Balance: 700,
	}

	mockRepo.On("UpdateWallet", ctx, expectedParams).Return(expectedWallet, nil)

	result, err := service.TopUpWalletBalance(ctx, walletID, amount)

	assert.NoError(t, err)
	assert.Equal(t, expectedWallet.ID, result.ID)
	assert.Equal(t, expectedWallet.Balance, result.Balance)

	mockRepo.AssertExpectations(t)
}
