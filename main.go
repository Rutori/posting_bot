package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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
	if err != nil {
		log.Fatal(err)
	}

	consumers.Run()
	_, _ = fmt.Scanln()
}
