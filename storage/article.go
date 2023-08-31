package storage

import (
	"TelegoBot/models"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"time"
)

func NewArticleStorage(db *sqlx.DB) *ArticleStorage {
	return &ArticleStorage{
		db: db,
	}
}

func (s *ArticleStorage) SaveArticle(ctx context.Context, article models.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO articles (title, url, summary, source_id, published_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING;`,
		article.Title,
		article.Url,
		article.Summary,
		article.SourceId,
		article.PublishedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *ArticleStorage) GetUnpostedArticles(ctx context.Context, since time.Time, limit uint64) ([]models.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	var articles []ArticleDB

	if err := conn.SelectContext(
		ctx,
		&articles,
		`
		SELECT 
				a.id AS a_id, 
				s.priority AS s_priority,
				s.id AS s_id,
				a.title AS a_title,
				a.url AS a_url,
				a.summary AS a_summary,
				a.published_at AS a_published_at,
				a.posted_at AS a_posted_at,
				a.created_at AS a_created_at
			FROM articles a JOIN sources s ON s.id = a.source_id
			WHERE a.posted_at IS NULL 
				AND a.published_at >= $1::timestamp
			ORDER BY a.created_at DESC, s_priority DESC LIMIT $2;
		`,
		since.UTC().Format(time.RFC3339),
		limit,
	); err != nil {
		return nil, err
	}

	return lo.Map(articles, func(article ArticleDB, _ int) models.Article {
		return models.Article{
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

func (s *ArticleStorage) MarkArticleAsPosted(ctx context.Context, article models.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`UPDATE articles SET posted_at = $1::timestamp WHERE id = $2`,
		time.Now().UTC().Format(time.RFC3339),
		article.Id,
	); err != nil {
		return err
	}

	return nil
}
