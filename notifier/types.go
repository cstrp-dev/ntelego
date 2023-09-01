package notifier

import (
	"TelegoBot/models"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type Notifier struct {
	articles   ArticleProvider
	summarizer Summarizer
	api        *tgbotapi.BotAPI
	apiKey     string
	prompt     string
	chanId     int64
	interval   time.Duration
	lookupTime time.Duration
}

type ArticleProvider interface {
	GetUnpostedArticles(ctx context.Context, since time.Time, limit uint64) ([]models.Article, error)
	MarkArticleAsPosted(ctx context.Context, article models.Article) error
}

type Summarizer interface {
	GetData(key, prompt, text string) (string, error)
}
