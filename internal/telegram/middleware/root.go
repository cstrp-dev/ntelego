package middleware

import (
	"TelegoBot/internal/storage"
	"TelegoBot/internal/telegram"
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func Root(userStorage *storage.UserStorage, next telegram.Callback) telegram.Callback {
	return func(ctx context.Context, b *tgbotapi.BotAPI, u tgbotapi.Update) error {
		chatId := u.FromChat().ID
		if err := userStorage.AddUser(ctx, chatId); err != nil {
			logrus.Errorf("Failed to add user %d: %v", chatId, err)
		}

		return next(ctx, b, u)
	}
}
