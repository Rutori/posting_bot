package posts

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"postingbot/config"
	"strings"
	"time"
)

func Schedule(message *tgbotapi.Message) (err error) {
	var lastDate int
	err = config.DB.Select(&lastDate, `SELECT max(date) FROM queue`)
	if lastDate == 0 {
		lastDate = time.Now().Second()
	}

	var photos []string
	for _, photo := range *message.Photo {
		photos = append(photos, photo.FileID)
	}

	_, err = config.DB.Exec(`INSERT INTO queue(date, text, photo, video, gif) VALUES (?,?,?,?,?)`,
		lastDate+config.JSON.Schedule.Interval,
		message.Text,
		strings.Join(photos, ","),
		message.Video.FileID,
		message.Animation.FileID,
	)

	return err
}
