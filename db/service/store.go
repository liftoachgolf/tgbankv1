package service

import (
	"context"
	"database/sql"
	"tgBank/db/postgres"
	"tgBank/models"
)

type StoreServiceImpl struct {
	repo postgres.StoreSQL
}

func NewStoreServiceImpl(repo postgres.StoreSQL) *StoreServiceImpl {
	return &StoreServiceImpl{
		repo: repo,
	}
}

func (s *StoreServiceImpl) TransferTx(ctx context.Context, arg models.TransferTxParams) (models.TransferTxResult, error) {
	return s.TransferTx(ctx, arg)
}

func (s *StoreServiceImpl) execTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	return s.execTx(ctx, fn)
}
