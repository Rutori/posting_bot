package consumers

import (
	"postingbot/telegram"
	"time"
)

func Run() {
	withInterval(telegram.FetchUpdates, time.Second)
}

func withInterval(launch func(), interval time.Duration) {
	for {
		launch()
		time.Sleep(interval)
	}
}
