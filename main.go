package main

import (
	"TelegoBot/config"
	"TelegoBot/helper"
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

	s, _ := helper.New()
	result, err := s.GetData(cfg.OpenAiApiKey, "I luv u!!")
	if err != nil {
		logrus.Errorf("Failed to get data from API. %v", err)
	}
	logrus.Print(result)
}
