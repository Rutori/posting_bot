package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mattn/go-sqlite3"
	"postingbot/config"
	"postingbot/consumers/posts"
	"strings"
)

func SelectCommand(update tgbotapi.Update) {
	if !checkAuth(update.Message.From.ID) {
		return
	}

	successMsg := parseCommand(update)
	successMsg.ReplyToMessageID = update.Message.MessageID
	successMsg.ChatID = update.Message.Chat.ID
	_, err := config.Bot.Send(successMsg)
	if err != nil {
		fmt.Printf("\ncontroller.SelectCommand #1: %s\n", err.Error())
		return
	}

	return
}

func checkAuth(id int) bool {
	for _, manager := range config.JSON.Managers {
		if manager == id {
			return true
		}
	}

	return false
}

func parseCommand(update tgbotapi.Update) (message tgbotapi.MessageConfig) {
	var err error
	var fallBackMsg = tgbotapi.MessageConfig{
		Text: "Sorry, server error",
	}

	if len(update.Message.Text) == 0 || update.Message.Text[:1] != "/" {
		message, err = queueMeme(update)
		if err != nil {
			return fallBackMsg
		}

		return message
	}

	firstSpace := strings.IndexRune(update.Message.Text, ' ')
	if firstSpace < 0 {
		firstSpace = len(update.Message.Text)
	}

	switch update.Message.Text[1:firstSpace] {
	case "listqueue":
		if message, err = GetQueue(); err != nil {
			return fallBackMsg
		}
	}

	return message
}

func queueMeme(update tgbotapi.Update) (successMsg tgbotapi.MessageConfig, err error) {
	date, err := posts.Schedule(update.Message)
	switch {
	case err == nil:

	case err.(sqlite3.Error).ExtendedCode == 0:

	case err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique:
		return successMsg, err

	default:
		fmt.Printf("\ncontroller.SelectCommand #2: %s\n", err.Error())
		return successMsg, err
	}
	successMsg = tgbotapi.NewMessage(int64(update.Message.From.ID), fmt.Sprintf("пост добавлен на %s", date))
	successMsg.ReplyToMessageID = update.Message.MessageID

	return successMsg, nil
}
