package storage

import (
	"TelegoBot/internal/models"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

func NewSourceStorage(db *sqlx.DB) *SourceStorage {
	return &SourceStorage{
		db: db,
	}
}

func (s *SourceStorage) GetAllSources(ctx context.Context) ([]models.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	var sources []SourceDB
	if err := conn.SelectContext(ctx, &sources, `SELECT * FROM sources`); err != nil {
		return nil, err
	}

	return lo.Map(sources, func(source SourceDB, _ int) models.Source {
		return models.Source(source)
	}), nil
}

func (s *SourceStorage) GetSourceById(ctx context.Context, id int64) (*models.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	var source SourceDB
	if err := conn.GetContext(ctx, &source, `SELECT * FROM sources WHERE id = $1`, id); err != nil {
		return nil, err
	}

	return (*models.Source)(&source), nil
}

func (s *SourceStorage) AddSource(ctx context.Context, source models.Source) (int64, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, err
	}

	defer conn.Close()

	var id int64

	row := conn.QueryRowContext(
		ctx,
		`INSERT INTO sources (name, feed_url, priority) VALUES ($1,$2,$3) RETURNING id`,
		source.Name, source.FeedUrl, source.Priority,
	)

	if err := row.Err(); err != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SourceStorage) SetSourcePriority(ctx context.Context, id int64, priority int) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `UPDATE sources SET priority = $1 WHERE id = $2`, priority, id); err != nil {
		return err
	}

	return nil
}

func (s *SourceStorage) DeleteSource(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `DELETE from sources WHERE id = $1`, id); err != nil {
		return err
	}

	return nil
}
