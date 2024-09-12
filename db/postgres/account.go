package postgres

import (
	"context"
	"database/sql"
	"tgBank/models"
)

type AccountSQLImpl struct {
	db *sql.DB
}

func NewAccountSQL(db *sql.DB) *AccountSQLImpl {
	return &AccountSQLImpl{
		db: db,
	}
}

func (q *AccountSQLImpl) AddAccountBalance(ctx context.Context, arg models.AddAccountBalanceParams) (models.Account, error) {
	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id= $2
		RETURNING id, owner, balance, currency, created_at
	`
	row := q.db.QueryRowContext(ctx, query, arg.Amount, arg.ID)
	var i models.Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

func (q *AccountSQLImpl) CreateAccount(ctx context.Context, arg models.CreateAccountParams) (models.Account, error) {
	query := `
		INSERT INTO accounts (
			owner,
			balance,
			currency
		) VALUES (
			$1, $2, $3
		) RETURNING id, owner, balance, currency, created_at
	`
	row := q.db.QueryRowContext(ctx, query, arg.Owner, arg.Balance, arg.Currency)
	var i models.Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

func (q *AccountSQLImpl) DeleteAccount(ctx context.Context, id int64) error {
	query := `
		DELETE FROM accounts
		WHERE id = $1
	`
	_, err := q.db.ExecContext(ctx, query, id)
	return err
}

func (q *AccountSQLImpl) GetAccount(ctx context.Context, id int64) (models.Account, error) {
	query := `
		SELECT id, owner, balance, currency, created_at 
		FROM accounts
		WHERE id = $1 LIMIT 1
	`
	row := q.db.QueryRowContext(ctx, query, id)
	var i models.Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

func (q *AccountSQLImpl) GetAccountForUpdate(ctx context.Context, id int64) (models.Account, error) {
	query := `
		SELECT id, owner, balance, currency, created_at 
		FROM accounts
		WHERE id = $1 LIMIT 1
		FOR NO KEY UPDATE
	`
	row := q.db.QueryRowContext(ctx, query, id)
	var i models.Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

func (q *AccountSQLImpl) ListAccounts(ctx context.Context, arg models.ListAccountsParams) ([]models.Account, error) {
	query := `
		SELECT id, owner, balance, currency, created_at 
		FROM accounts
		ORDER BY id LIMIT $1 OFFSET $2
	`
	rows, err := q.db.QueryContext(ctx, query, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.Account{}
	for rows.Next() {
		var i models.Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q *AccountSQLImpl) UpdateAccount(ctx context.Context, arg models.UpdateAccountParams) (models.Account, error) {
	query := `
		UPDATE accounts
		SET balance = $2
		WHERE id = $1
		RETURNING id, owner, balance, currency, created_at
	`
	row := q.db.QueryRowContext(ctx, query, arg.ID, arg.Balance)
	var i models.Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

func (q *AccountSQLImpl) addMoney(ctx context.Context, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 models.Account, account2 models.Account, err error) {
	account1, err = q.AddAccountBalance(ctx, models.AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, models.AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return

}
