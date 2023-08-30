package telegram

import (
	"TelegoBot/telegram"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func Rm(r RemoveSource) telegram.Cb {
	return func(ctx context.Context, api *tgbotapi.BotAPI, upd tgbotapi.Update) error {
		id, err := strconv.ParseInt(upd.Message.CommandArguments(), 10, 64)
		if err != nil {
			return err
		}

		if err := r.Rm(ctx, id); err != nil {
			return nil
		}

		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Source has been removed successfully")
		if _, err := api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
