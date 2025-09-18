package storage

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type ArticleStorage struct {
	db *sqlx.DB
}

type SourceStorage struct {
	db *sqlx.DB
}

type UserStorage struct {
	db *sqlx.DB
}

type SourceDB struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedUrl   string    `db:"feed_url"`
	Priority  int       `db:"priority"`
	CreatedAt time.Time `db:"created_at"`
}

type ArticleDB struct {
	Id             int64          `db:"a_id"`
	SourceId       int64          `db:"s_id"`
	SourcePriority int64          `db:"s_priority"`
	Title          string         `db:"a_title"`
	Url            string         `db:"a_url"`
	Summary        sql.NullString `db:"a_summary"`
	PublishedAt    time.Time      `db:"a_published_at"`
	PostedAt       sql.NullTime   `db:"a_posted_at"`
	CreatedAt      time.Time      `db:"a_created_at"`
}

type UserDB struct {
	ChatId    int64     `db:"chat_id"`
	CreatedAt time.Time `db:"created_at"`
}
