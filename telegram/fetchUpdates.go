package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"postingbot/config"
)

func FetchUpdates() {
	updates, err := config.Bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Limit:   100,
		Timeout: 5,
	})
	if err != nil {
		fmt.Printf("\ntelegram.FetchUpdates #1: %s\n", err.Error())
		return
	}
	for update := range updates {
		SelectCommand(update)
	}

}
