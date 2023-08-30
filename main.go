package main

import (
	"TelegoBot/config"
	"TelegoBot/storage"
	"TelegoBot/telegram"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Error("Failed to load .env file!")
	}
}

func main() {
	cfg := config.New()
	api, err := tgbotapi.NewBotAPI(cfg.TelegramApiKey)
	api.Debug = true

	if err != nil {
		logrus.Error("Failed to create bot api:", err)
		return
	}

	db, err := sqlx.Connect("postgres", cfg.DatabaseUrl)
	if err != nil {
		logrus.Error("Failed to connect to db:", err)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			logrus.Error("Failed to close db:", err)
		}
	}()

	var (
		as = storage.NewArticlesStorage(db)
		ss = storage.NewSourcesStorage(db)
	)

	bot := telegram.New(api)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer func() {
		cancel()
		logrus.Info("Bot stopped")
	}()

	if err := bot.Init(ctx); err != nil {
		logrus.Error("Failed to init bot:", err)
	}

	logrus.Info(as, ss)
}
