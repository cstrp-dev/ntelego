package notifier

import (
	"TelegoBot/internal/models"
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Notifier struct {
	articles    ArticleProvider
	summarizer  Summarizer
	b           *tgbotapi.BotAPI
	userStorage UserStorage
	interval    time.Duration
	lookupTime  time.Duration
}

type Summarizer interface {
	Summarize(text string) (string, error)
}

type ArticleProvider interface {
	GetUnpostedArticles(ctx context.Context, since time.Time, limit uint64) ([]models.Article, error)
	MarkArticleAsPosted(ctx context.Context, article models.Article) error
}

type UserStorage interface {
	GetAllUsers(ctx context.Context) ([]int64, error)
}
