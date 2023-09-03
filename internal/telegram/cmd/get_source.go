package telegram

import (
	"TelegoBot/internal/telegram"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func GetSource(p SourceProvider) telegram.Callback {
	return func(ctx context.Context, b *tgbotapi.BotAPI, u tgbotapi.Update) error {
		id, err := strconv.ParseInt(u.Message.CommandArguments(), 10, 64)
		if err != nil {
			return err
		}

		source, err := p.GetSourceById(ctx, id)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, Format(*source))
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := b.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
