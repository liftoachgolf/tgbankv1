package postgres

import (
	"context"
	"database/sql"
	"tgBank/models"
)

type EntrySQLImpl struct {
	db *sql.DB
}

func NewEntrySQL(db *sql.DB) *EntrySQLImpl {
	return &EntrySQLImpl{
		db: db,
	}
}

func (q *EntrySQLImpl) CreateEntry(ctx context.Context, arg models.CreateEntryParams) (models.Entry, error) {
	query := `INSERT INTO entries (account_id, amount) VALUES ($1, $2) RETURNING id, account_id, amount, created_at`
	row := q.db.QueryRowContext(ctx, query, arg.AccountID, arg.Amount)
	var i models.Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

func (q *EntrySQLImpl) GetEntry(ctx context.Context, id int64) (models.Entry, error) {
	query := `SELECT id, account_id, amount, created_at FROM entries
WHERE id = $1 LIMIT 1`
	row := q.db.QueryRowContext(ctx, query, id)
	var i models.Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

func (q *EntrySQLImpl) ListEntries(ctx context.Context, arg models.ListEntriesParams) ([]models.Entry, error) {
	query := `SELECT id, account_id, amount, created_at FROM entries
WHERE account_id = $1
ORDER BY id
    LIMIT $2
OFFSET $3`
	rows, err := q.db.QueryContext(ctx, query, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.Entry{}
	for rows.Next() {
		var i models.Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
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
