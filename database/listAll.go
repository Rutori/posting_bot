package database

import (
	"postingbot/config"
	"postingbot/entities"
)

func ListAll() (list []*entities.Message, err error) {
	list = make([]*entities.Message, 0)
	err = config.DB.Select(&list, `
		SELECT 
		       date, 
		       coalesce(text,'') as text, 
		       coalesce(photo,'') as photo, 
		       coalesce(video,'') as video, 
		       coalesce(gif,'') as gif,
		       id 
		FROM queue WHERE date > strftime('%s', 'now');`)

	return list, err
}
