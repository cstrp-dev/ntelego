package telegram

import (
	"TelegoBot/internal/models"
	"TelegoBot/internal/telegram"
	"context"
	"fmt"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
)

func SourceLs(l SourceList) telegram.Callback {
	return func(ctx context.Context, b *tgbotapi.BotAPI, u tgbotapi.Update) error {
		s, err := l.GetAllSources(ctx)
		if err != nil {
			return err
		}

		sort.SliceStable(s, func(i, j int) bool {
			return s[i].Priority > s[j].Priority
		})

		info := lo.Map(s, func(src models.Source, _ int) string {
			return Format(src)
		})
		msg := tgbotapi.NewMessage(
			u.Message.Chat.ID,
			fmt.Sprintf(
				"List of sources (%d):\n\n%s",
				len(s),
				strings.Join(info, "\n\n"),
			),
		)

		if _, err := b.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
