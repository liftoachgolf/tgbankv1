package service

import (
	"context"
	"tgBank/models"
	"tgBank/pkg/postgres"
)

type AccountServiceImpl struct {
	repo postgres.AccountSQL
}

func NewAccountServiceImpl(repo postgres.AccountSQL) *AccountServiceImpl {
	return &AccountServiceImpl{
		repo: repo,
	}
}

func (s *AccountServiceImpl) CreateAccount(ctx context.Context, arg models.CreateAccountParams) (models.Account, error) {
	return s.repo.CreateAccount(ctx, arg)
}
func (s *AccountServiceImpl) DeleteAccount(ctx context.Context, id int64) error {
	return s.repo.DeleteAccount(ctx, id)
}
func (s *AccountServiceImpl) GetAccount(ctx context.Context, id int64) (models.Account, error) {
	return s.repo.GetAccount(ctx, id)
}
func (s *AccountServiceImpl) GetAccountForUpdate(ctx context.Context, id int64) (models.Account, error) {
	return s.repo.GetAccountForUpdate(ctx, id)
}
func (s *AccountServiceImpl) ListAccounts(ctx context.Context, arg models.ListAccountsParams) ([]models.Account, error) {
	return s.repo.ListAccounts(ctx, arg)
}
func (s *AccountServiceImpl) UpdateAccount(ctx context.Context, arg models.UpdateAccountParams) (models.Account, error) {
	return s.repo.UpdateAccount(ctx, arg)
}
