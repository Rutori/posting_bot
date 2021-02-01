package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"postingbot/constants"
	"postingbot/consumers/posts"
	"postingbot/database"
	"time"
)

const lineTemplate = "Пост с %s в очереди на %s\n"
const totalTemplate = "\nВсего %d постов"

var textTypes = map[int]string{
	constants.MessageTypes.Text:  "текстом",
	constants.MessageTypes.Pic:   "фото",
	constants.MessageTypes.Video: "видно",
	constants.MessageTypes.Gif:   "гифкой",
}

func GetQueue() (message tgbotapi.MessageConfig, err error) {
	list, err := database.ListAll()
	if err != nil {
		return message, err
	}

	for _, post := range list {
		message.Text += fmt.Sprintf(lineTemplate, textTypes[posts.GetType(post)], post.Date.In(time.Local).Format("15:04 02/01/2006 "))
	}

	message.Text += fmt.Sprintf(totalTemplate, len(list))

	return message, err
}
