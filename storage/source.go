package storage

import (
	"TelegoBot/model"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func NewSourcesStorage(db *sqlx.DB) *SourceStorage {
	return &SourceStorage{
		db: db,
	}
}

func (s *SourceStorage) Get(ctx context.Context) ([]model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Panic(err)
		}
	}()

	var sources []SourceDB
	if err := conn.SelectContext(ctx, &sources, `SELECT * FROM sources;`); err != nil {
		return nil, err
	}

	return nil, err
}

func (s *SourceStorage) GetByID(ctx context.Context, id int64) (*model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Panic(err)
		}
	}()

	var source SourceDB
	if err := conn.GetContext(ctx, &source, `SELECT * FROM sources WHERE id = $1;`, id); err != nil {
		return nil, err
	}

	return (*model.Source)(&source), nil
}

func (s *SourceStorage) Add(ctx context.Context, source model.Source) (int64, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Panic(err)
		}
	}()

	var id int64
	row := conn.QueryRowContext(ctx, `INSERT INTO sources(name, feed_url, priority) VALUES ($1, $2, $3) RETURNING id;`, source.Name, source.FeedURL, source.Priority)

	if err := row.Err(); err != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, err
}

func (s *SourceStorage) Rm(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Panic(err)
		}
	}()

	if _, err := conn.ExecContext(ctx, `DELETE FROM sources WHERE id = $1;`, id); err != nil {
		return err
	}

	return nil
}

func (s *SourceStorage) Set(ctx context.Context, id int64, priority int) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Panic(err)
		}
	}()

	if _, err := conn.ExecContext(ctx, `UPDATE sources SET priority = $1 WHERE id = $2;`, priority, id); err != nil {
		return err
	}

	return nil
}
