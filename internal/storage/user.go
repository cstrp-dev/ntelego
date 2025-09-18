package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) AddUser(ctx context.Context, chatId int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx, `
		INSERT INTO users (chat_id, created_at)
		VALUES ($1, NOW())
		ON CONFLICT (chat_id) DO NOTHING`,
		chatId)
	return err
}

func (s *UserStorage) GetAllUsers(ctx context.Context) ([]int64, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var users []UserDB
	err = conn.SelectContext(ctx, &users, `SELECT chat_id FROM users`)
	if err != nil {
		return nil, err
	}

	chatIds := make([]int64, len(users))
	for i, user := range users {
		chatIds[i] = user.ChatId
	}
	return chatIds, nil
}
