package posts

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"postingbot/config"
	"strconv"
	"time"
)

func Schedule(message *tgbotapi.Message) (date string, err error) {
	var lastDate int64
	err = config.DB.QueryRow(`SELECT coalesce(max(date),0) FROM queue`).Scan(&lastDate)
	if err != nil {
		return "", err
	}
	var postDate int64
	if lastDate < time.Now().Unix() {
		postDate = time.Now().Unix()
	} else {
		postDate = lastDate
	}

	t := time.Unix(postDate, 0)
	dayStart := time.Date(t.Year(), t.Month(), t.Day(), config.JSON.Schedule.From, 0, 0, t.Nanosecond(), t.Location())
	dayEnd := time.Date(t.Year(), t.Month(), t.Day(), config.JSON.Schedule.To, 0, 0, t.Nanosecond(), t.Location())

	postDate += config.JSON.Schedule.Interval

	if postDate < dayStart.Unix() {
		postDate = dayStart.Unix()
	}

	if postDate > dayEnd.Unix() {
		postDate = dayStart.Add(24 * time.Hour).Unix()
	}

	var args = []interface{}{
		postDate,
		fmt.Sprintf("%s|%s", strconv.Itoa(message.MessageID), strconv.Itoa(message.From.ID)),
	}
	switch {
	case message.Photo != nil && len(*message.Photo) > 0:
		photos := *message.Photo
		args = append(args, message.Text, photos[0].FileID)
		_, err = config.DB.Exec(`INSERT INTO queue(date,id, text, photo) VALUES (?,?,?,?)`, args...)

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

	return time.Unix(postDate, 0).String(), err
}
