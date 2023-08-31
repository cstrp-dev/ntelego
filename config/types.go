package config

import "time"

type Config struct {
	TelegramApiKey    string
	TelegramChannelID int64
	DatabaseUrl       string
	Keywords          []string

	OpenAiApiKey         string
	Prompt               string
	FetchInterval        time.Duration
	NotificationInterval time.Duration
}
