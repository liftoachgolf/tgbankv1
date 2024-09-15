package postgrestg

import (
	"context"
	"database/sql"
	"tgBank/models"
)

type entryTgSQL struct {
	db *sql.DB
}

func NewEntrySQL(db *sql.DB) EntryTgSQL {
	return &entryTgSQL{
		db: db,
	}
}

func (q *entryTgSQL) CreateEntry(ctx context.Context, arg models.CreateEntryTgParams) (models.EntryTg, error) {
	query := `INSERT INTO entries (chat_id, amount) VALUES ($1, $2) RETURNING id, chat_id, amount, created_at`
	row := q.db.QueryRowContext(ctx, query, arg.ChatId, arg.Amount)
	var i models.EntryTg
	err := row.Scan(
		&i.ID,
		&i.ChatId,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

func (q *entryTgSQL) GetEntry(ctx context.Context, id int64) (models.EntryTg, error) {
	query := `SELECT id, chat_id, amount, created_at FROM entries
WHERE id = $1 LIMIT 1`
	row := q.db.QueryRowContext(ctx, query, id)
	var i models.EntryTg
	err := row.Scan(
		&i.ID,
		&i.ChatId,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

func (q *entryTgSQL) ListEntries(ctx context.Context, arg models.ListEntriesTgParams) ([]models.EntryTg, error) {
	query := `SELECT id, chat_id, amount, created_at FROM entries
WHERE chat_id = $1
ORDER BY id
    LIMIT $2
OFFSET $3`
	rows, err := q.db.QueryContext(ctx, query, arg.ChatId, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.EntryTg{}
	for rows.Next() {
		var i models.EntryTg
		if err := rows.Scan(
			&i.ID,
			&i.ChatId,
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
