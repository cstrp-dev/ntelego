package fetcher

import (
	"TelegoBot/helpers"
	"TelegoBot/models"
	src "TelegoBot/source"
	"context"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

func New(as ArticleStorage, sp SourceProvider, i time.Duration, kw []string) *Fetcher {
	return &Fetcher{
		articles: as,
		sources:  sp,
		interval: i,
		keywords: kw,
	}
}

func (f *Fetcher) Init(ctx context.Context) error {
	sources, err := f.sources.GetAllSources(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, source := range sources {
		wg.Add(1)

		go func(source Source) {
			defer wg.Done()

			items, err := source.Fetch(ctx)
			if err != nil {
				logrus.Errorf("Failed to fetch %s: %s", source.Name(), err)
				return
			}

			if err := f.processItems(ctx, source, items); err != nil {
				logrus.Errorf("Failed to process items for %s: %s", source.Name(), err)
				return
			}

		}(src.NewRSSource(source))

	}

	wg.Wait()

	return nil
}

func (f *Fetcher) processItems(ctx context.Context, source Source, items []models.Item) error {
	for _, item := range items {
		item.Date = item.Date.UTC()

		if f.shouldSkip(item) {
			logrus.Infof("Item %s %q should be skipped (%q)", item.Title, item.Url, source.Name())
			continue
		}

		if err := f.articles.SaveArticle(ctx, models.Article{
			SourceId:    source.Id(),
			Title:       item.Title,
			Url:         item.Url,
			Summary:     item.Summary,
			PublishedAt: item.Date,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fetcher) shouldSkip(item models.Item) bool {
	set := helpers.NewSet(item.Categories)

	for _, key := range f.keywords {
		if set.Contains(key) || strings.Contains(strings.ToLower(item.Title), key) {
			return true
		}
	}

	return false
}
