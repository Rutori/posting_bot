package config

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

var JSON struct {
	DatabaseFile string   `json:"databaseFile"`
	Key          string   `json:"key"`
	Channel      int64    `json:"channel"`
	Managers     []int    `json:"managers"`
	Schedule     Schedule `json:"schedule"`
}

var Bot *tgbotapi.BotAPI

var DB *sqlx.DB

type Schedule struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Interval int    `json:"interval"`
}
