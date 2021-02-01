package posts

import (
	"postingbot/constants"
	"postingbot/entities"
)

func GetType(message *entities.Message) int {
	switch {
	case message.Photo != "":
		return constants.MessageTypes.Pic

	case message.Video != "":
		return constants.MessageTypes.Video

	case message.Gif != "":
		return constants.MessageTypes.Gif

	default:
		return constants.MessageTypes.Text
	}
}
