package telegram

import (
	"TelegoBot/internal/models"
	"fmt"
)

func Format(s models.Source) string {
	return fmt.Sprintf(
		"🛠 %s\nID: %d\nFeed URL: %s\nPriority: %d",
		s.Name,
		s.Id,
		s.FeedUrl,
		s.Priority,
	)
}
