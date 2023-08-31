package models

import "time"

type Source struct {
	Id        int64
	Name      string
	FeedUrl   string
	Priority  int
	CreatedAt time.Time
}
