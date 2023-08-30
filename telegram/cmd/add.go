package telegram

import (
	"TelegoBot/model"
	"TelegoBot/telegram"
	"TelegoBot/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func Add(s SourceStorage) telegram.Cb {
	type Add struct {
		Name     string `json:"name"`
		Url      string `json:"url"`
		Priority int    `json:"priority"`
	}

	return func(ctx context.Context, api *tgbotapi.BotAPI, upd tgbotapi.Update) error {
		args, err := utils.JsonParse[Add](upd.Message.CommandArguments())
		if err != nil {
			return err
		}

		source := model.Source{
			Name:     args.Name,
			FeedURL:  args.Url,
			Priority: args.Priority,
		}
		sourceID, err := s.Add(ctx, source)
		if err != nil {
			logrus.Error("Failed to add new source")
			return err
		}

		var (
			msg   = fmt.Sprintf("Add new source with id: %d", sourceID)
			reply = tgbotapi.NewMessage(upd.Message.Chat.ID, msg)
		)

		reply.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := api.Send(reply); err != nil {
			logrus.Error("Failed to send reply")
			return err
		}

		return nil
	}
}
