package middleware

import (
	"TelegoBot/telegram"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Root(channelId int64, next telegram.Callback) telegram.Callback {
	return func(ctx context.Context, b *tgbotapi.BotAPI, u tgbotapi.Update) error {
		amdns, err := b.GetChatAdministrators(
			tgbotapi.ChatAdministratorsConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: channelId,
				},
			},
		)
		if err != nil {
			return err
		}

		for _, adm := range amdns {
			if adm.User.ID == u.SentFrom().ID {
				return next(ctx, b, u)
			}
		}

		msg := tgbotapi.NewMessage(u.FromChat().ID, "Access denied!")

		if _, err := b.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
