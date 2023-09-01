package fetcher

import (
	"TelegoBot/models"
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
	SaveArticle(ctx context.Context, article models.Article) error
}

type SourceProvider interface {
	GetAllSources(ctx context.Context) ([]models.Source, error)
}

type Source interface {
	Id() int64
	Name() string
	Fetch(ctx context.Context) ([]models.Item, error)
}
