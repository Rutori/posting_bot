package posts

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/otiai10/gosseract"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
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

	postDate = int64(
		math.Floor(
			float64(
				(postDate+config.JSON.Schedule.Interval)/config.JSON.Schedule.Interval),
		) * float64(config.JSON.Schedule.Interval))

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
		imageText, err := addTextFromImage(photos[0].FileID)
		if err != nil {
			return "", err
		}

		message.Text += imageText
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

func addTextFromImage(fileID string) (string, error) {
	url, err := config.Bot.GetFileDirectURL(fileID)
	if err != nil {
		return "", err
	}

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		er := response.Body.Close()
		if er != nil {
			fmt.Println(er)
		}
	}()

	img, err := jpeg.Decode(response.Body)
	if err != nil {
		return "", err
	}

	grayImg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, grayImg)
	if err != nil {
		return "", err
	}

	client := gosseract.NewClient()
	err = client.SetWhitelist("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	if err != nil {
		return "", err
	}

	err = client.SetPageSegMode(gosseract.PSM_SPARSE_TEXT)
	if err != nil {
		return "", err
	}

	err = client.SetImageFromBytes(buf.Bytes())
	if err != nil {
		return "", err
	}

	str, err := client.Text()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("[%s]", str), nil
}
