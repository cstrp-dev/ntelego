package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"time"
)

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{api: api}
}

func (b *Bot) Init(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(context.Background(), 5*time.Minute)
			b.handleUpdates(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) RegistryCmd(cmd string, callback Callback) {
	if b.cmds == nil {
		b.cmds = make(map[string]Callback)
	}

	b.cmds[cmd] = callback
}

func (b *Bot) handleUpdates(ctx context.Context, upd tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			logrus.Errorf("Panic in handle updates - recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()

	if (upd.Message == nil || !upd.Message.IsCommand()) && upd.CallbackQuery == nil {
		return
	}

	if !upd.Message.IsCommand() {
		return
	}

	var cb Callback
	cmd := upd.Message.Command()
	cmdView, ok := b.cmds[cmd]
	if !ok {
		return
	}

	cb = cmdView
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Internal error !")

	if err := cb(ctx, b.api, upd); err != nil {
		logrus.Errorf("Failed to execute cmd view: %v", err)

		if _, err := b.api.Send(msg); err != nil {
			logrus.Errorf("Failed to send message: %v", err)
		}
	}

}
