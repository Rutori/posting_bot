package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

		err = posts.Schedule(update.Message)
		if err != nil {
			fmt.Printf("\ntelegram.FetchUpdates #2: %s\n", err.Error())
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
