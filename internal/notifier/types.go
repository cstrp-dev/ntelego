package notifier

import (
	"TelegoBot/internal/models"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type Notifier struct {
	articles   ArticleProvider
	summarizer Summarizer
	b          *tgbotapi.BotAPI
	channelId  int64
	interval   time.Duration
	lookupTime time.Duration
}

type Summarizer interface {
	Summarize(text string) (string, error)
}

type ArticleProvider interface {
	GetUnpostedArticles(ctx context.Context, since time.Time, limit uint64) ([]models.Article, error)
	MarkArticleAsPosted(ctx context.Context, article models.Article) error
}
