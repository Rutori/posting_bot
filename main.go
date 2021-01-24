package main

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"postingbot/config"
	"postingbot/consumers"
)

func main() {
	conf, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(conf, &config.JSON)
	if err != nil {
		log.Fatal(err)
	}

	config.Bot, err = tgbotapi.NewBotAPI(config.JSON.Key)
	if err != nil {
		log.Fatal(err)
	}

	config.DB, err = sqlx.Connect("sqlite3", config.JSON.DatabaseFile)

	consumers.Run()
}