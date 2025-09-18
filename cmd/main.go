package main

import (
	"TelegoBot/cmd/config"
	"TelegoBot/internal/fetcher"
	"TelegoBot/internal/helpers"
	"TelegoBot/internal/notifier"
	storage2 "TelegoBot/internal/storage"
	telegram3 "TelegoBot/internal/telegram"
	telegram2 "TelegoBot/internal/telegram/cmd"
	"TelegoBot/internal/telegram/middleware"
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Errorf("Failed to load .env file. %v", err)
	}
}

func main() {
	cfg := config.New()
	api, err := tgbotapi.NewBotAPI(cfg.TelegramApiKey)
	if err != nil {
		logrus.Errorf("Failed to create bot . %v", err)
		return
	}

	db, err := sqlx.Connect("postgres", cfg.DatabaseUrl)
	if err != nil {
		logrus.Errorf("Failed to connect to database. %v", err)
		return
	}

	h, err := helpers.New(cfg.OpenAiApiKey, cfg.Prompt)
	if err != nil {
		logrus.Errorf("Failed to create helpers. %v", err)
	}

	var (
		articleStorage = storage2.NewArticleStorage(db)
		sourceStorage  = storage2.NewSourceStorage(db)
		userStorage    = storage2.NewUserStorage(db)
		fetcher        = fetcher.New(articleStorage, sourceStorage, cfg.FetchInterval, cfg.Keywords)
		notifier       = notifier.New(
			articleStorage,
			h, api,
			userStorage,
			cfg.NotificationInterval,
			2*cfg.FetchInterval,
		)
	)

	newBot := telegram3.New(api)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func(ctx context.Context) {
		if err := fetcher.Init(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				logrus.Errorf("Failed to init fetcher. %v", err)
				return
			}
			logrus.Info("Fetcher stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := notifier.Init(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				logrus.Errorf("Failed to init notifier. %v", err)
				return
			}
			logrus.Info("Notifier stopped")
		}
	}(ctx)

	setMyCommands(newBot, sourceStorage, userStorage)
	if err := newBot.Init(ctx); err != nil {
		logrus.Errorf("Failed to init bot. %v", err)
		return
	}

}

func setMyCommands(b *telegram3.Bot, storage *storage2.SourceStorage, userStorage *storage2.UserStorage) {
	b.RegistryCmd("add", middleware.Root(userStorage, telegram2.AddSource(storage)))
	b.RegistryCmd("get", middleware.Root(userStorage, telegram2.GetSource(storage)))
	b.RegistryCmd("set", middleware.Root(userStorage, telegram2.SetPriority(storage)))
	b.RegistryCmd("ls", middleware.Root(userStorage, telegram2.SourceLs(storage)))
	b.RegistryCmd("rm", middleware.Root(userStorage, telegram2.DeleteSource(storage)))
}
