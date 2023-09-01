package telegram

import (
	"TelegoBot/models"
	"context"
)

type SourceStorage interface {
	AddSource(ctx context.Context, source models.Source) (int64, error)
}

type SourceProvider interface {
	GetSourceById(ctx context.Context, id int64) (*models.Source, error)
}

type SourceList interface {
	GetAllSources(ctx context.Context) ([]models.Source, error)
}

type SourceRemover interface {
	DeleteSource(ctx context.Context, sourceID int64) error
}

type PrioritySetter interface {
	SetSourcePriority(ctx context.Context, sourceID int64, priority int) error
}
