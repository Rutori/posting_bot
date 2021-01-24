package consumers

import (
	"postingbot/telegram"
	"time"
)

func Run() {
	go withInterval(telegram.FetchUpdates, time.Second)
	go withInterval(telegram.PostMeme, time.Second)
}

func withInterval(launch func(), interval time.Duration) {
	for {
		launch()
		time.Sleep(interval)
	}
}
