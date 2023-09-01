package notifier

import (
	"TelegoBot/helpers"
	"TelegoBot/models"
	"context"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/go-shiori/go-readability"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"
)

func New(
	provider ArticleProvider,
	summarizer Summarizer,
	b *tgbotapi.BotAPI,
	channelId int64,
	interval time.Duration,
	lookupTime time.Duration,
) *Notifier {
	return &Notifier{
		articles:   provider,
		summarizer: summarizer,
		b:          b,
		channelId:  channelId,
		interval:   interval,
		lookupTime: lookupTime,
	}
}

func (n *Notifier) Init(ctx context.Context) error {
	ticker := time.NewTicker(n.interval)
	defer ticker.Stop()

	if err := n.selectArticle(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := n.selectArticle(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (n *Notifier) selectArticle(ctx context.Context) error {
	articles, err := n.articles.GetUnpostedArticles(ctx, time.Now().Add(-n.lookupTime), 4)
	if err != nil {
		return err
	}

	if len(articles) == 0 {
		logrus.Info("No articles to notify")
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
	var reader io.Reader

	if article.Summary != "" {
		reader = strings.NewReader(article.Summary)
	} else {
		res, err := http.Get(article.Url)
		if err != nil {
			return "", err
		}

		reader = res.Body
	}

	doc, err := readability.FromReader(reader, nil)
	if err != nil {
		return "", err
	}

	summary, err := n.summarizer.Summarize(helpers.CleanUpText(doc.TextContent))
	if err != nil {
		return "", err
	}

	return "\n\n" + summary, nil
}

func (n *Notifier) send(article models.Article, summary string) error {
	msg := tgbotapi.NewMessage(
		n.channelId,
		fmt.Sprintf(
			"*%s*%s\n\n%s",
			helpers.Escape(article.Title),
			helpers.Escape(summary),
			helpers.Escape(article.Url),
		),
	)

	msg.ParseMode = "MarkdownV2"

	if _, err := n.b.Send(msg); err != nil {
		return err
	}

	return nil
}
