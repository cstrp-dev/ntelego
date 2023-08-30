package telegram

import (
	"TelegoBot/model"
	"context"
)

type SourceList interface {
	Get(ctx context.Context) ([]model.Source, error)
}

type SourceStorage interface {
	Add(ctx context.Context, source model.Source) (int64, error)
}

type SourceProvider interface {
	GetByID(ctx context.Context, id int64) (*model.Source, error)
}

type RemoveSource interface {
	Rm(ctx context.Context, id int64) error
}

type SetPriority interface {
	Set(ctx context.Context, id int64, priority int) error
}
