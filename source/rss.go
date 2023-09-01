package source

import (
	"TelegoBot/models"
	"context"
	"github.com/SlyMarbo/rss"
	"github.com/samber/lo"
	"strings"
)

func NewRSSource(m models.Source) RSSource {
	return RSSource{
		SourceId:   m.Id,
		SourceName: m.Name,
		Url:        m.FeedUrl,
	}
}

func (s RSSource) Fetch(ctx context.Context) ([]models.Item, error) {
	feed, err := s.load(ctx, s.Url)
	if err != nil {
		return nil, err
	}

	return lo.Map(feed.Items, func(item *rss.Item, _ int) models.Item {
			return models.Item{
				Title:      item.Title,
				Url:        item.Link,
				Date:       item.Date,
				Categories: item.Categories,
				SourceName: s.SourceName,
				Summary:    strings.TrimSpace(item.Summary),
			}
		}),
		nil
}

func (s RSSource) Id() int64 {
	return s.SourceId
}

func (s RSSource) Name() string {
	return s.SourceName
}

func (s RSSource) load(ctx context.Context, url string) (*rss.Feed, error) {
	var (
		feedChan = make(chan *rss.Feed)
		errChan  = make(chan error)
	)

	go func() {
		rssFeed, err := rss.Fetch(url)
		if err != nil {
			errChan <- err
			return
		}
		feedChan <- rssFeed
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errChan:
		return nil, err
	case feed := <-feedChan:
		return feed, nil
	}

}
