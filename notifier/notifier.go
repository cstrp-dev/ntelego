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
	"regexp"
	"strings"
	"time"
)

func New(
	ap ArticleProvider,
	s Summarizer,
	api *tgbotapi.BotAPI,
	k, p string,
	id int64,
	i time.Duration,
	l time.Duration,
) *Notifier {
	return &Notifier{
		articles:   ap,
		summarizer: s,
		api:        api,
		apiKey:     k,
		prompt:     p,
		chanId:     id,
		interval:   i,
		lookupTime: l,
	}
}

func (n *Notifier) Init(ctx context.Context) error {
	ticker := time.NewTicker(n.interval)
	defer ticker.Stop()

	if err := n.GetAndSend(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := n.GetAndSend(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

}

func (n *Notifier) GetAndSend(ctx context.Context) error {
	articles, err := n.articles.GetUnpostedArticles(ctx, time.Now().Add(-n.lookupTime), 1)
	if err != nil {
		return err
	}

	if len(articles) == 0 {
		return nil
	}

	article := articles[0]

	summary, err := n.extract(n.apiKey, n.prompt, article)
	if err != nil {
		logrus.Errorf("Error extracting summary: %v", err)
	}

	if err := n.sendArticle(article, summary); err != nil {
		logrus.Errorf("Error sending article: %v", err)
		return err
	}

	return n.articles.MarkArticleAsPosted(ctx, article)
}

func (n *Notifier) extract(key, prompt string, article models.Article) (string, error) {
	var reader io.Reader

	if article.Summary != "" {
		reader = strings.NewReader(article.Summary)
	} else {
		resp, err := http.Get(article.Url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		reader = resp.Body
	}

	doc, err := readability.FromReader(reader, nil)
	if err != nil {
		return "", err
	}

	summary, err := n.summarizer.GetData(key, prompt, cleanUp(doc.TextContent))
	if err != nil {
		return "", err
	}

	return "\n\n" + summary, nil
}

func cleanUp(text string) string {
	redundantNewLines := regexp.MustCompile(`\n{3,}`)
	return redundantNewLines.ReplaceAllString(text, "\n")
}

func (n *Notifier) sendArticle(article models.Article, summary string) error {
	const format = "*%s*%s\n\n%s"

	msg := tgbotapi.NewMessage(
		n.chanId,
		fmt.Sprintf(format,
			helpers.Escape(article.Title),
			helpers.Escape(summary),
			helpers.Escape(article.Url),
		))
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	if _, err := n.api.Send(msg); err != nil {
		return err
	}

	return nil
}
