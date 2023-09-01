package telegram

import (
	"TelegoBot/helpers"
	"TelegoBot/models"
	"TelegoBot/telegram"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddSource(s SourceStorage) telegram.Callback {
	type Args struct {
		Name     string `json:"name"`
		Url      string `json:"url"`
		Priority int    `json:"priority"`
	}

	return func(ctx context.Context, b *tgbotapi.BotAPI, u tgbotapi.Update) error {
		args, err := helpers.JSONParse[Args](u.Message.CommandArguments())
		if err != nil {
			return err
		}

		source := models.Source{
			Name:     args.Name,
			FeedUrl:  args.Url,
			Priority: args.Priority,
		}

		id, err := s.AddSource(ctx, source)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(
			u.Message.Chat.ID,
			fmt.Sprintf("Source has been added successfully. ID: %d", id),
		)

		msg.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := b.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
