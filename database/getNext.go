package database

import (
	"postingbot/config"
	"postingbot/entities"
	"time"
)

func GetNext() (*entities.Message, error) {
	msgStruct := &entities.Message{}
	timePost := time.Now().Unix()
	err := config.DB.QueryRow("SELECT `text`,`photo`,`video`,`gif` FROM `queue` WHERE `date` > ? AND `date` < ?",
		timePost-config.JSON.Schedule.Interval,
		timePost+config.JSON.Schedule.Interval).Scan(&msgStruct.Text, &msgStruct.Photo, &msgStruct.Video, msgStruct.Gif)
	if err != nil {
		return nil, err
	}

	return msgStruct, nil
}
