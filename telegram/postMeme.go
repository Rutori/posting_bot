package telegram

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"postingbot/config"
	"postingbot/consumers/posts"
	"time"
)

func PostMeme() {
	message := &posts.Message{}
	err := config.DB.QueryRow(`
		SELECT 
		       coalesce(text,''), 
		       coalesce(photo,''), 
		       coalesce(video,''), 
		       coalesce(gif,'') 
		from queue 
		where date = ?`, time.Now().Unix()).Scan(
		&message.Text,
		&message.Photo,
		&message.Video,
		&message.Gif,
	)
	switch err {
	case nil:

	case sql.ErrNoRows:
		return

	default:
		fmt.Printf("\ntelegram.PostMeme #1: %s\n", err.Error())
		return

	}

	switch {
	case len(message.Photo) > 0:
		err = postPhoto(message.Photo)
	case len(message.Video) > 0:
		err = postVideo(message.Video)
	case len(message.Gif) > 0:
		err = postGif(message.Gif)

	default:
		_, err = config.Bot.Send(tgbotapi.NewMessage(config.JSON.Channel, message.Text))
	}
	if err != nil {
		fmt.Printf("\ntelegram.PostMeme #2: %s\n", err.Error())
	}
}

func postPhoto(fileID string) error {
	_, err := config.Bot.Send(tgbotapi.NewPhotoShare(config.JSON.Channel, fileID))
	return err
}

func postVideo(fileID string) error {
	_, err := config.Bot.Send(tgbotapi.NewVideoShare(config.JSON.Channel, fileID))
	return err
}

func postGif(fileID string) error {
	_, err := config.Bot.Send(tgbotapi.NewAnimationShare(config.JSON.Channel, fileID))
	return err
}
