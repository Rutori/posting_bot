package posts

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"postingbot/config"
	"strconv"
	"strings"
	"time"
)

func Schedule(message *tgbotapi.Message) (err error) {
	var lastDate int64
	err = config.DB.QueryRow(`SELECT coalesce(max(date),0) FROM queue`).Scan(&lastDate)
	if err != nil {
		return err
	}

	if lastDate < time.Now().Unix() {
		lastDate = time.Now().Unix()
	}
	var args = []interface{}{
		lastDate + config.JSON.Schedule.Interval,
		fmt.Sprintf("%s|%s", strconv.Itoa(message.MessageID), strconv.Itoa(message.From.ID)),
	}
	switch {
	case len(*message.Photo) > 0:
		var photos []string
		for _, photo := range *message.Photo {
			photos = append(photos, photo.FileID)
		}
		args = append(args, strings.Join(photos, ","))
		_, err = config.DB.Exec(`INSERT INTO queue(date,id, photo) VALUES (?,?,?)`, args...)

	case message.Video != nil:
		args = append(args, message.Video.FileID)
		_, err = config.DB.Exec(`INSERT INTO queue(date,id, video) VALUES (?,?,?)`, args...)

	case message.Animation != nil:
		args = append(args, message.Animation.FileID)
		_, err = config.DB.Exec(`INSERT INTO queue(date,id, gif) VALUES (?,?,?)`, args...)

	default:
		args = append(args, message.Text)
		_, err = config.DB.Exec(`INSERT INTO queue(date,id, text) VALUES (?,?,?)`, args...)

	}

	return err
}
