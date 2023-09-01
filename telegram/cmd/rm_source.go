package telegram

import (
	"TelegoBot/telegram"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func DeleteSource(d SourceRemover) telegram.Callback {
	return func(ctx context.Context, b *tgbotapi.BotAPI, u tgbotapi.Update) error {
		id, err := strconv.ParseInt(u.Message.CommandArguments(), 10, 64)
		if err != nil {
			return err
		}

		if err := d.DeleteSource(ctx, id); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Source has been remove successfully.")
		if _, err := b.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
