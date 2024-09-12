package service

import (
	"context"
	"tgBank/db/postgres"
	"tgBank/models"
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
	return s.CreateAccount(ctx, arg)
}
func (s *AccountServiceImpl) DeleteAccount(ctx context.Context, id int64) error {
	return s.DeleteAccount(ctx, id)
}
func (s *AccountServiceImpl) GetAccount(ctx context.Context, id int64) (models.Account, error) {
	return s.GetAccount(ctx, id)
}
func (s *AccountServiceImpl) GetAccountForUpdate(ctx context.Context, id int64) (models.Account, error) {
	return s.GetAccountForUpdate(ctx, id)
}
func (s *AccountServiceImpl) ListAccounts(ctx context.Context, arg models.ListAccountsParams) ([]models.Account, error) {
	return s.ListAccounts(ctx, arg)
}
func (s *AccountServiceImpl) UpdateAccount(ctx context.Context, arg models.UpdateAccountParams) (models.Account, error) {
	return s.UpdateAccount(ctx, arg)
}
func (s *AccountServiceImpl) addMoney(ctx context.Context, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 models.Account, account2 models.Account, err error) {
	return s.addMoney(ctx, accountID1, amount1, accountID2, amount2)
}
