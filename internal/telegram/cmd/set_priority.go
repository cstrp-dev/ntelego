package telegram

import (
	"TelegoBot/internal/helpers"
	"TelegoBot/internal/telegram"
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SetPriority(p PrioritySetter) telegram.Callback {
	type Args struct {
		SourceId int64 `json:"source_id"`
		Priority int   `json:"priority"`
	}

	return func(ctx context.Context, b *tgbotapi.BotAPI, u tgbotapi.Update) error {
		argsStr := u.Message.CommandArguments()
		if argsStr == "" {
			msg := tgbotapi.NewMessage(
				u.Message.Chat.ID,
				"❌ Error: Please provide arguments in JSON format.\nExample: `/set {\"source_id\":1,\"priority\":5}`",
			)
			b.Send(msg)
			return nil
		}

		args, err := helpers.JSONParse[Args](argsStr)
		if err != nil {
			msg := tgbotapi.NewMessage(
				u.Message.Chat.ID,
				fmt.Sprintf("❌ Error parsing arguments: %v\nPlease use valid JSON format.\nExample: `/set {\"source_id\":1,\"priority\":5}`", err),
			)
			b.Send(msg)
			return nil
		}

		if err := p.SetSourcePriority(ctx, args.SourceId, args.Priority); err != nil {
			msg := tgbotapi.NewMessage(
				u.Message.Chat.ID,
				fmt.Sprintf("❌ Error setting priority: %v", err),
			)
			b.Send(msg)
			return nil
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "✅ Priority set successfully.")
		if _, err := b.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
