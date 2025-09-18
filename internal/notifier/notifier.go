package notifier

import (
	"TelegoBot/internal/helpers"
	"TelegoBot/internal/models"
	"context"
	"fmt"
	"time"

	http "github.com/bogdanfinn/fhttp"
	"github.com/go-shiori/go-readability"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func New(
	provider ArticleProvider,
	summarizer Summarizer,
	b *tgbotapi.BotAPI,
	userStorage UserStorage,
	interval time.Duration,
	lookupTime time.Duration,
) *Notifier {
	return &Notifier{
		articles:    provider,
		summarizer:  summarizer,
		b:           b,
		userStorage: userStorage,
		interval:    interval,
		lookupTime:  lookupTime,
	}
}

func (n *Notifier) Init(ctx context.Context) error {
	ticker := time.NewTicker(n.interval)
	defer ticker.Stop()

	if err := n.selectArticle(ctx); err != nil {
		logrus.Errorf("Error getting and sending articles: %v", err)
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := n.selectArticle(ctx); err != nil {
				logrus.Errorf("Error getting and sending articles: %v", err)
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (n *Notifier) selectArticle(ctx context.Context) error {
	format := time.Now().Add(-n.lookupTime).UTC().Format("2006-01-02 15:04:05")
	parsed, _ := time.Parse(time.RFC3339, format)
	articles, err := n.articles.GetUnpostedArticles(ctx, parsed, 1)
	if err != nil {
		return err
	}

	if len(articles) == 0 {
		logrus.Infof("%v articles to send.", len(articles))
		return nil
	}

	article := articles[0]

	summary, err := n.extract(article)
	if err != nil {
		logrus.Errorf("Failed to extract summary. %v", err)
	}

	if err := n.send(article, summary); err != nil {
		logrus.Errorf("Failed to send summary. %v", err)
		return err
	}

	return n.articles.MarkArticleAsPosted(ctx, article)
}

func (n *Notifier) extract(article models.Article) (string, error) {
	res, err := http.Get(article.Url)
	if err != nil {
		if article.Summary != "" {
			logrus.Warnf("Failed to fetch article content from URL, using RSS summary: %v", err)
			return n.summarizer.Summarize(helpers.CleanUpText(article.Summary))
		}
		return "", fmt.Errorf("failed to fetch article content: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		if article.Summary != "" {
			logrus.Warnf("Article URL returned status %d, using RSS summary", res.StatusCode)
			return n.summarizer.Summarize(helpers.CleanUpText(article.Summary))
		}
		return "", fmt.Errorf("article URL returned status %d", res.StatusCode)
	}

	doc, err := readability.FromReader(res.Body, nil)
	if err != nil {
		if article.Summary != "" {
			logrus.Warnf("Failed to extract content with readability, using RSS summary: %v", err)
			return n.summarizer.Summarize(helpers.CleanUpText(article.Summary))
		}
		return "", fmt.Errorf("failed to extract article content: %v", err)
	}

	content := helpers.CleanUpText(doc.TextContent)
	if len(content) < 100 {
		if article.Summary != "" && len(article.Summary) > len(content) {
			logrus.Warnf("Extracted content too short (%d chars), using RSS summary", len(content))
			return n.summarizer.Summarize(helpers.CleanUpText(article.Summary))
		}
	}

	summary, err := n.summarizer.Summarize(content)
	if err != nil {
		return "", fmt.Errorf("failed to summarize article: %v", err)
	}

	logrus.Infof("Article processed - Title: %s, Content length: %d chars, Summary length: %d chars",
		article.Title, len(content), len(summary))

	return "\n\n" + summary, nil
}

func (n *Notifier) send(article models.Article, summary string) error {
	users, err := n.userStorage.GetAllUsers(context.Background())
	if err != nil {
		return err
	}

	messageText := fmt.Sprintf(
		"📰 %s\n\n%s\n\n🔗 %s",
		article.Title,
		summary,
		article.Url,
	)

	for _, chatId := range users {
		msg := tgbotapi.NewMessage(
			chatId,
			messageText,
		)

		if _, err := n.b.Send(msg); err != nil {
			logrus.Errorf("Failed to send message to chat %d: %v", chatId, err)
		}
	}

	return nil
}
