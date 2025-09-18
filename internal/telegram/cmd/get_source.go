package telegram

import (
	"TelegoBot/internal/telegram"
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetSource(p SourceProvider) telegram.Callback {
	return func(ctx context.Context, b *tgbotapi.BotAPI, u tgbotapi.Update) error {
		argsStr := u.Message.CommandArguments()
		if argsStr == "" {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "❌ Error: Please provide a source ID.\nExample: `/get 1`")
			b.Send(msg)
			return nil
		}

		id, err := strconv.ParseInt(argsStr, 10, 64)
		if err != nil {
			msg := tgbotapi.NewMessage(
				u.Message.Chat.ID,
				fmt.Sprintf("❌ Error: Invalid ID format '%s'. Please provide a valid number.\nExample: `/get 1`", argsStr),
			)
			b.Send(msg)
			return nil
		}

		source, err := p.GetSourceById(ctx, id)
		if err != nil {
			msg := tgbotapi.NewMessage(
				u.Message.Chat.ID,
				fmt.Sprintf("❌ Error getting source: %v", err),
			)
			b.Send(msg)
			return nil
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, Format(*source))

		if _, err := b.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
