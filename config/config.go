package config

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	cfg  Config
	once sync.Once
)

const (
	DefaultStr           = ""
	DefaultInt           = 0
	TelegramApiKey       = "TELEGRAM_API_KEY"
	TelegramChannelID    = "TELEGRAM_CHANNEL_ID"
	DatabaseUrl          = "DATABASE_URL"
	FetchInterval        = "FETCH_INTERVAL"
	NotificationInterval = "NOTIFICATION_INTERVAL"
	Keywords             = "KEYWORDS"
	OpenaiApiKey         = "OPENAI_API_KEY"
)

func New() Config {
	once.Do(func() {
		cfg = Config{
			TelegramApiKey:       getString(TelegramApiKey, DefaultStr),
			TelegramChannelID:    getInt(TelegramChannelID, DefaultInt),
			DatabaseUrl:          getString(DatabaseUrl, DefaultStr),
			OpenAiApiKey:         getString(OpenaiApiKey, DefaultStr),
			FetchInterval:        getDuration(FetchInterval, time.Second*30),
			NotificationInterval: getDuration(NotificationInterval, time.Second*30),
			Keywords:             getStringSlice(Keywords, nil),
		}
	})

	return cfg
}

func getString(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return i
		}
	}
	return defaultValue
}

func getStringSlice(key string, defaultValue []string) []string {
	if value, ok := os.LookupEnv(key); ok {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
	}
	return defaultValue
}
