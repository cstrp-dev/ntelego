package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api  *tgbotapi.BotAPI
	cmds map[string]Cb
}

type Cb func(ctx context.Context, api *tgbotapi.BotAPI, upd tgbotapi.Update) error
