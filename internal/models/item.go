package models

import "time"

type Item struct {
	Title      string
	Url        string
	Summary    string
	Categories []string
	SourceName string
	Date       time.Time
}
