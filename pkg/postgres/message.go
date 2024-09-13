package postgres

import (
	"context"
	"database/sql"
	"tgBank/models"
)

type MessageSqlImpl struct {
	db *sql.DB
}

func NewMessageSql(db *sql.DB) *MessageSqlImpl {
	return &MessageSqlImpl{
		db: db,
	}
}

func (r *MessageSqlImpl) execTx(fn func(tx *sql.Tx) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Обязательно откатываем транзакцию, если функция возвращает ошибку
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *MessageSqlImpl) Create(ctx context.Context, arg models.Message) (int, error) {
	var chatExists bool

	err := r.execTx(func(tx *sql.Tx) error {
		checkChatQuery := "SELECT EXISTS (SELECT 1 FROM chats WHERE chat_id = $1)"
		err := tx.QueryRowContext(ctx, checkChatQuery, arg.ChatId).Scan(&chatExists)
		if err != nil {
			return err
		}

		if !chatExists {
			createChatQuery := "INSERT INTO chats(chat_id) VALUES ($1) RETURNING id"
			_, err := tx.ExecContext(ctx, createChatQuery, arg.ChatId)
			if err != nil {
				return err
			}
		}
		createMessageQuery := "INSERT INTO messages (chat_id, message_id, text) VALUES ($1, $2, $3)"
		_, err = tx.ExecContext(ctx, createMessageQuery, arg.ChatId, arg.MessageId, arg.Text)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return 0, nil
}
