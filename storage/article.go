package storage

import (
	"TelegoBot/model"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"time"
)

func NewArticlesStorage(db *sqlx.DB) *ArticleStorage {
	return &ArticleStorage{
		db: db,
	}
}

func (s *ArticleStorage) Add(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Panic(err)
		}
	}()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO articles (source_id, title, url, summary, published_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING;`,
		article.SourceId,
		article.Title,
		article.Url,
		article.Summary,
		article.PublishedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *ArticleStorage) IsNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Panic(err)
		}
	}()

	var articles []ArticleDB

	if err := conn.SelectContext(
		ctx,
		&articles,
		`SELECT 
				a.id as a_id, 
				s.priority as s_priority,
				s.id as s_id,
				a.title as a_title,
				a.url as a_url,
				a.summary as a_summary,
				a.published_at as a_published_at,
				a.posted_at as a_posted_at,
				a.created_at as a_created_at
			FROM articles a JOIN sources s ON s.id = a.source_id
			WHERE a.posted_at IS NULL 
				AND a.published_at >= $1::timestamp
			ORDER BY a.created_at DESC, s_priority DESC LIMIT $2;`,
		since.UTC().Format(time.RFC3339),
		limit,
	); err != nil {
		return nil, err
	}

	return lo.Map(articles, func(article ArticleDB, _ int) model.Article {
		return model.Article{
			Id:          article.Id,
			Title:       article.Title,
			Url:         article.Url,
			Summary:     article.Summary.String,
			SourceId:    article.SourceId,
			CreatedAt:   article.CreatedAt,
			PublishedAt: article.PublishedAt,
		}
	}), nil
}

func (s *ArticleStorage) IsPosted(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Panic(err)
		}
	}()

	if _, err := conn.ExecContext(
		ctx,
		`UPDATE articles SET posted_at = $1::timestamp WHERE id = $2;`,
		time.Now().UTC().Format(time.RFC3339),
		article.Id,
	); err != nil {
		return err
	}

	return nil
}
