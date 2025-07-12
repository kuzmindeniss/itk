package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kuzmindeniss/itk/internal/db/repository"
)

type WalletRepositoryInterface interface {
	GetWalletByID(ctx context.Context, id uuid.UUID) (repository.Wallet, error)
	UpdateWallet(ctx context.Context, arg repository.UpdateWalletParams) (repository.Wallet, error)
}

type WalletServiceInterface interface {
	GetWalletByID(ctx context.Context, id uuid.UUID) (repository.Wallet, error)
	TopUpWalletBalance(ctx context.Context, id uuid.UUID, amount int32) (repository.Wallet, error)
}

type WalletService struct {
	repo WalletRepositoryInterface
}

func NewWalletService(repo WalletRepositoryInterface) *WalletService {
	return &WalletService{
		repo: repo,
	}
}

func (s *WalletService) GetWalletByID(ctx context.Context, id uuid.UUID) (repository.Wallet, error) {
	return s.repo.GetWalletByID(ctx, id)
}

func (s *WalletService) TopUpWalletBalance(ctx context.Context, id uuid.UUID, amount int32) (repository.Wallet, error) {
	return s.repo.UpdateWallet(ctx, repository.UpdateWalletParams{
		ID:     id,
		Amount: amount,
	})
}
