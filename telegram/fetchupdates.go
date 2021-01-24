package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mattn/go-sqlite3"
	"postingbot/config"
	"postingbot/consumers/posts"
)

func FetchUpdates() {
	updates, err := config.Bot.GetUpdates(tgbotapi.UpdateConfig{
		Offset:  0,
		Limit:   100,
		Timeout: 5,
	})
	if err != nil {
		fmt.Printf("\ntelegram.FetchUpdates #1: %s\n", err.Error())
		return
	}
	for _, update := range updates {
		if !checkAuth(update.Message.From.ID) {
			continue
		}
		date, err := posts.Schedule(update.Message)
		switch {
		case err == nil:

		case err.(sqlite3.Error).ExtendedCode == 0:

		case err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintUnique:
			continue

		default:
			fmt.Printf("\ntelegram.FetchUpdates #2: %s\n", err.Error())
			return
		}
		successMsg := tgbotapi.NewMessage(int64(update.Message.From.ID), fmt.Sprintf("пост добавлен на %s", date))
		successMsg.ReplyToMessageID = update.Message.MessageID
		_, err = config.Bot.Send(successMsg)
		if err != nil {
			fmt.Printf("\ntelegram.FetchUpdates #3: %s\n", err.Error())
			return
		}
	}

}

func checkAuth(id int) bool {
	for _, manager := range config.JSON.Managers {
		if manager == id {
			return true
		}
	}
	return false
}
