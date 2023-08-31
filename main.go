package main

import (
	"TelegoBot/config"
	"TelegoBot/helpers"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Errorf("Failed to load .env file. %v", err)
	}
}

func main() {
	cfg := config.New()

	s, _ := helpers.New()
	result, err := s.GetData(cfg.OpenAiApiKey, cfg.Prompt, "")
	if err != nil {
		logrus.Errorf("Failed to get data from API. %v", err)
	}

	fmt.Print(result)
}
