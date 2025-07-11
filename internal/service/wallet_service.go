package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kuzmindeniss/itk/internal/db/repository"
)

type WalletService struct {
	repo *repository.Queries
}

func NewWalletService(repo *repository.Queries) *WalletService {
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
