package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"postingbot/config"
	"postingbot/consumers/posts"
)

func PostMeme(message *posts.Message) (err error) {
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

	return err
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
