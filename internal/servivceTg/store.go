package serviceTg

import (
	"context"
	"database/sql"
	"tgBank/models"
	postgrestg "tgBank/pkg/postgresTg"
)

type storeTgService struct {
	repo         postgrestg.MessageTgSQL
	repoStore    postgrestg.StoreTgSQL
	repoAccount  postgrestg.AccountTgSQL
	repoTransfer postgrestg.TransferTgSQL
	repoEntry    postgrestg.EntryTgSQL
}

func NewStoreTgService(repo postgrestg.MessageTgSQL, repoStore postgrestg.StoreTgSQL, repoAccount postgrestg.AccountTgSQL, repoTransfer postgrestg.TransferTgSQL, repoEntry postgrestg.EntryTgSQL) StoreTgService {
	return &storeTgService{
		repo:         repo,
		repoStore:    repoStore,
		repoAccount:  repoAccount,
		repoTransfer: repoTransfer,
		repoEntry:    repoEntry,
	}
}

func (store *storeTgService) TransferTx(ctx context.Context, arg models.TransferTgTxParams) (models.TransferTgTxResult, error) {
	var result models.TransferTgTxResult

	err := store.repoStore.ExecTx(ctx, func(tx *sql.Tx) error {
		var err error

		result.Transfer, err = store.repoTransfer.CreateTransfer(ctx, models.CreateTransferTgParams{
			FromChatId: arg.FromChatId,
			ToChatId:   arg.ToChatId,
			Amount:     arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = store.repoEntry.CreateEntry(ctx, models.CreateEntryTgParams{
			ChatId: arg.FromChatId,
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = store.repoEntry.CreateEntry(ctx, models.CreateEntryTgParams{
			ChatId: arg.ToChatId,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromChatId < arg.ToChatId {
			result.FromChatId, result.ToChatId, err = store.repoAccount.AddMoney(ctx, arg.FromChatId, -arg.Amount, arg.ToChatId, arg.Amount)
		} else {
			result.ToChatId, result.FromChatId, err = store.repoAccount.AddMoney(ctx, arg.ToChatId, arg.Amount, arg.FromChatId, -arg.Amount)
		}
		return nil
	})

	return result, err
}
