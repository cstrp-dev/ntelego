package storage

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type ArticleStorage struct {
	db *sqlx.DB
}

type SourceStorage struct {
	db *sqlx.DB
}

type ArticleDB struct {
	Id             int64          `db:"a_id"`
	SourcePriority int64          `db:"s_priority"`
	Title          string         `db:"a_title"`
	Url            string         `db:"a_url"`
	Summary        sql.NullString `db:"a_summary"`
	SourceId       int64          `db:"s_id"`
	CreatedAt      time.Time      `db:"a_created_at"`
	PublishedAt    time.Time      `db:"a_published_at"`
	PostedAt       time.Time      `db:"a_posted_at"`
}

type SourceDB struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	Priority  int       `db:"priority"`
	CreatedAt time.Time `db:"created_at"`
}
