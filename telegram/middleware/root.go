package telegram

import (
	"TelegoBot/telegram"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func OnlyRoot(channelID int64, cb telegram.Cb) telegram.Cb {
	return func(ctx context.Context, api *tgbotapi.BotAPI, upd tgbotapi.Update) error {
		admins, err := api.GetChatAdministrators(
			tgbotapi.ChatAdministratorsConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: channelID,
				},
			},
		)
		if err != nil {
			return err
		}

		for _, admin := range admins {
			if admin.User.ID == upd.SentFrom().ID {
				return cb(ctx, api, upd)
			}
		}

		msg := tgbotapi.NewMessage(upd.FromChat().ID, "You don't have privileges for user this cmd.")
		if _, err := api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
