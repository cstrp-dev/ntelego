package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api  *tgbotapi.BotAPI
	cmds map[string]Callback
}

type Callback func(ctx context.Context, b *tgbotapi.BotAPI, u tgbotapi.Update) error
