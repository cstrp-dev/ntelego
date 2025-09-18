package telegram

import (
	"TelegoBot/internal/telegram"
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DeleteSource(d SourceRemover) telegram.Callback {
	return func(ctx context.Context, b *tgbotapi.BotAPI, u tgbotapi.Update) error {
		argsStr := u.Message.CommandArguments()
		if argsStr == "" {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "❌ Error: Please provide a source ID.\nExample: `/rm 1`")
			b.Send(msg)
			return nil
		}

		id, err := strconv.ParseInt(argsStr, 10, 64)
		if err != nil {
			msg := tgbotapi.NewMessage(
				u.Message.Chat.ID,
				fmt.Sprintf("❌ Error: Invalid ID format '%s'. Please provide a valid number.\nExample: `/rm 1`", argsStr),
			)
			b.Send(msg)
			return nil
		}

		if err := d.DeleteSource(ctx, id); err != nil {
			msg := tgbotapi.NewMessage(
				u.Message.Chat.ID,
				fmt.Sprintf("❌ Error deleting source: %v", err),
			)
			b.Send(msg)
			return nil
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "✅ Source has been removed successfully.")
		if _, err := b.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
