package models

import "time"

type Article struct {
	Id          int64
	Title       string
	Url         string
	Summary     string
	SourceId    int64
	CreatedAt   time.Time
	PublishedAt time.Time
	PostedAt    time.Time
}
