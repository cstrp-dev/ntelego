package model

import "time"

type Article struct {
	Id       int64
	SourceId int64
	Title    string
	Url      string
	Summary  string

	CreatedAt   time.Time
	PublishedAt time.Time
	PostedAt    time.Time
}
