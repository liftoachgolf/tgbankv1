package serviceTg

import (
	"context"
	"database/sql"
	"tgBank/models"
	postgrestg "tgBank/pkg/postgresTg"
)

type messageTgService struct {
	repo        postgrestg.MessageTgSQL
	repoStore   postgrestg.StoreTgSQL
	repoAccount postgrestg.AccountTgSQL
}

func NewMessageTgService(repo postgrestg.MessageTgSQL, repoStore postgrestg.StoreTgSQL, repoAccount postgrestg.AccountTgSQL) MessageTgService {
	return &messageTgService{
		repo:        repo,
		repoStore:   repoStore,
		repoAccount: repoAccount,
	}
}

func (m *messageTgService) CreateMessage(msg models.Message) error {
	err := m.repoStore.ExecTx(context.Background(), func(tx *sql.Tx) error {

		isChatExists, err := m.repoAccount.IsExitsts(context.Background(), int64(msg.ChatId))
		if err != nil {
			return err
		}
		if !isChatExists {
			err := m.repoAccount.CreateAccount(context.Background(), int64(msg.ChatId), msg.Username, 777, "USD")
			if err != nil {
				return err
			}
		}
		err = m.repo.CreateMessage(msg)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
