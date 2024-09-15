package postgrestg

import (
	"context"
	"database/sql"
	"fmt"
	"tgBank/models"
)

type accountTgSQL struct {
	db *sql.DB
}

func NewAccountTGSQL(db *sql.DB) AccountTgSQL {
	return &accountTgSQL{
		db: db,
	}
}

func (r *accountTgSQL) CreateAccount(ctx context.Context, id int64, username string, balance int64, currency string) error {
	query := `INSERT INTO accounts (chat_id, username, balance, currency) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, id, username, balance, currency)
	if err != nil {
		return fmt.Errorf("could not create account: %w", err)
	}
	return nil
}

func (r *accountTgSQL) IsExitsts(ctx context.Context, chatID int64) (bool, error) {
	var exists bool

	query := "SELECT  EXISTS (SELECT 1 FROM accounts WHERE chat_id=$1)"
	err := r.db.QueryRowContext(ctx, query, chatID).Scan(&exists)

	return exists, err
}

func (r *accountTgSQL) GetAccount(ctx context.Context, chatID int64) (*models.AccountTg, error) {
	query := `SELECT id, chat_id, username, balance, currency FROM accounts WHERE chat_id = $1`
	row := r.db.QueryRowContext(ctx, query, chatID)

	var account models.AccountTg
	err := row.Scan(&account.ID, &account.ChatID, &account.Username, &account.Balance, &account.Currency)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("could not get account: %w", err)
	}

	return &account, nil
}

func (r *accountTgSQL) AddAccountBalance(ctx context.Context, arg models.AddAccountTgBalanceParams) (models.AccountTg, error) {
	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE chat_id = $2
		RETURNING chat_id, username, balance, currency
	`
	row := r.db.QueryRowContext(ctx, query, arg.Amount, arg.ChatID)
	var i models.AccountTg
	err := row.Scan(
		&i.ChatID,
		&i.Username,
		&i.Balance,
		&i.Currency,
	)
	return i, err
}

func (r *accountTgSQL) AddMoney(ctx context.Context, chatId1 int64, amount1 int64, chatId2 int64, amount2 int64) (account1 models.AccountTg, account2 models.AccountTg, err error) {
	account1, err = r.AddAccountBalance(ctx, models.AddAccountTgBalanceParams{
		ChatID: chatId1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = r.AddAccountBalance(ctx, models.AddAccountTgBalanceParams{
		ChatID: chatId2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return

}
