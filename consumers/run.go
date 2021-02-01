package consumers

import (
	"postingbot/telegram"
	"time"
)

func Run() {
	go withInterval(telegram.FetchUpdates, time.Second)
	go withInterval(telegram.PostMeme, time.Minute)
}

func withInterval(launch func(), interval time.Duration) {
	for {
		launch()
		time.Sleep(interval)
	}
}
