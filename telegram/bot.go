package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"time"
)

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		api: api,
	}
}

func (b *Bot) Init(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-updates:
			upCtx, upCancel := context.WithTimeout(context.Background(), 5*time.Minute)
			b.handleUpdate(upCtx, update)
			upCancel()
		}
	}
}

func (b *Bot) RegistryCmd(cmd string, cb Cb) {
	if b.cmds == nil {
		b.cmds = make(map[string]Cb)
	}

	b.cmds[cmd] = cb
}

func (b *Bot) handleUpdate(ctx context.Context, upd tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			logrus.Infof("Panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()

	var cb Cb
	if (upd.Message == nil || !upd.Message.IsCommand()) && upd.CallbackQuery == nil {
		return
	}

	cmd := upd.Message.Command()
	view, ok := b.cmds[cmd]
	if !ok {
		return
	}

	cb = view
	if err := cb(ctx, b.api, upd); err != nil {
		logrus.Error("Failed to exec cmd.")
	}
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Internal error.")
	if _, err := b.api.Send(msg); err != nil {
		logrus.Error("Failed to send message.")
	}

}
