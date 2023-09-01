package telegram

import (
	"TelegoBot/helpers"
	"TelegoBot/models"
	"fmt"
)

func Format(s models.Source) string {
	return fmt.Sprintf(
		"🛠 *%s*\nID: `%d`\nFeed url: %s\nPriority: %d",
		helpers.Escape(s.Name),
		s.Id,
		helpers.Escape(s.FeedUrl),
		s.Priority,
	)
}
