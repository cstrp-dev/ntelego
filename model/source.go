package model

import "time"

type Source struct {
	Id        int64
	Name      string
	FeedURL   string
	Priority  int
	CreatedAt time.Time
}
