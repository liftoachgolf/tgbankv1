package serviceTg

import (
	"context"
	"tgBank/models"
	postgrestg "tgBank/pkg/postgresTg"
)

type accountTgService struct {
	repo postgrestg.AccountTgSQL
}

func NewAccountTgSerivce(repo postgrestg.AccountTgSQL) AccountTgService {
	return &accountTgService{
		repo: repo,
	}
}

func (r *accountTgService) IsExitsts(ctx context.Context, chatID int64) (bool, error) {
	return r.repo.IsExitsts(ctx, chatID)
}

func (r *accountTgService) CreateAccount(ctx context.Context, id int64, username string, balance int64, currency string) error {
	return r.repo.CreateAccount(ctx, id, username, balance, currency)
}

func (r *accountTgService) GetAccount(ctx context.Context, chatID int64) (*models.AccountTg, error) {
	return r.repo.GetAccount(ctx, chatID)
}

func (r *accountTgService) AddMoney(ctx context.Context, chatId1 int64, amount1 int64, chatId2 int64, amount2 int64) (account1 models.AccountTg, account2 models.AccountTg, err error) {
	return r.repo.AddMoney(ctx, chatId1, amount1, chatId2, amount2)
}
func (r *accountTgService) AddAccountBalance(ctx context.Context, arg models.AddAccountTgBalanceParams) (models.AccountTg, error) {
	return r.repo.AddAccountBalance(ctx, arg)
}
