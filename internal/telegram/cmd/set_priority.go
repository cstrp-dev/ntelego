package telegram

import (
	"TelegoBot/internal/helpers"
	"TelegoBot/internal/telegram"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SetPriority(p PrioritySetter) telegram.Callback {
	type Args struct {
		SourceId int64 `json:"source_id"`
		Priority int   `json:"priority"`
	}

	return func(ctx context.Context, b *tgbotapi.BotAPI, u tgbotapi.Update) error {
		args, err := helpers.JSONParse[Args](u.Message.CommandArguments())
		if err != nil {
			return err
		}

		if err := p.SetSourcePriority(ctx, args.SourceId, args.Priority); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Set priority successfully.")
		if _, err := b.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
