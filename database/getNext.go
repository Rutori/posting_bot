package database

import (
	"postingbot/config"
	"postingbot/consumers/posts"
	"time"
)

func GetNext() (*posts.Message, error) {
	msgStruct := &posts.Message{}
	err := config.DB.QueryRow("SELECT `text`,`photo`,`video`,`gif` FROM `queue` WHERE `date` > ? AND `date` < ?",
		time.Now().Second()-config.JSON.Schedule.Interval,
		time.Now().Second()+config.JSON.Schedule.Interval).Scan(&msgStruct.Text, &msgStruct.Photo, &msgStruct.Video, msgStruct.Gif)
	if err != nil {
		return nil, err
	}

	return msgStruct, nil
}
