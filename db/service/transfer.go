package service

import (
	"context"
	"tgBank/db/postgres"
	"tgBank/models"
)

type TransferServiceImpl struct {
	repo postgres.TransferSQL
}

func NewTransferServiceImpl(repo postgres.TransferSQL) *TransferServiceImpl {
	return &TransferServiceImpl{
		repo: repo,
	}
}

func (s *TransferServiceImpl) ListTransfers(ctx context.Context, arg models.ListTransfersParams) ([]models.Transfer, error) {
	return s.ListTransfers(ctx, arg)
}
func (s *TransferServiceImpl) CreateTransfer(ctx context.Context, arg models.CreateTransferParams) (models.Transfer, error) {
	return s.CreateTransfer(ctx, arg)
}
func (s *TransferServiceImpl) GetTransfer(ctx context.Context, id int64) (models.Transfer, error) {
	return s.GetTransfer(ctx, id)
}
