package entities

import "time"

type Message struct {
	Date  time.Time `db:"date"`
	Text  string    `db:"text"`
	Photo string    `db:"photo"`
	Video string    `db:"video"`
	Gif   string    `db:"gif"`
	ID    string    `db:"id"`
}
