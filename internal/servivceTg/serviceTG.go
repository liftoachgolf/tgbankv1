package serviceTg

import (
	"context"
	"tgBank/models"
	postgrestg "tgBank/pkg/postgresTg"
)

type AccountTgService interface {
	CreateAccount(ctx context.Context, id int64, username string, balance int64, currency string) error
	GetAccount(ctx context.Context, chatID int64) (*models.AccountTg, error)
	IsExitsts(ctx context.Context, chatID int64) (bool, error)
	AddMoney(ctx context.Context, chatId1 int64, amount1 int64, chatId2 int64, amount2 int64) (account1 models.AccountTg, account2 models.AccountTg, err error)
	AddAccountBalance(ctx context.Context, arg models.AddAccountTgBalanceParams) (models.AccountTg, error)
}

type MessageTgService interface {
	CreateMessage(msg models.Message) error
}

type StoreTgService interface {
	TransferTx(ctx context.Context, arg models.TransferTgTxParams) (models.TransferTgTxResult, error)
}

type ServiceTg struct {
	AccountTgService
	MessageTgService
	StoreTgService
}

func NewServiceTg(repo *postgrestg.Repository) *ServiceTg {
	return &ServiceTg{
		AccountTgService: NewAccountTgSerivce(repo.AccountTgSQL),
		MessageTgService: NewMessageTgService(repo.MessageTgSQL, repo.StoreTgSQL, repo.AccountTgSQL),
		StoreTgService:   NewStoreTgService(repo.MessageTgSQL, repo.StoreTgSQL, repo.AccountTgSQL, repo.TransferTgSQL, repo.EntryTgSQL),
	}
}
