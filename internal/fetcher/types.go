package fetcher

import (
	models2 "TelegoBot/internal/models"
	"context"
	"time"
)

type Fetcher struct {
	articles ArticleStorage
	sources  SourceProvider
	interval time.Duration
	keywords []string
}

type ArticleStorage interface {
	SaveArticle(ctx context.Context, article models2.Article) error
}

type SourceProvider interface {
	GetAllSources(ctx context.Context) ([]models2.Source, error)
}

type Source interface {
	Id() int64
	Name() string
	Fetch(ctx context.Context) ([]models2.Item, error)
}
