package constants

const (
	messageTypeText  = 1
	messageTypePic   = 2
	messageTypeVideo = 3
	messageTypeGif   = 4
)

type types struct {
	Text  int
	Pic   int
	Video int
	Gif   int
}

var MessageTypes = types{
	Text:  messageTypeText,
	Pic:   messageTypePic,
	Video: messageTypeVideo,
	Gif:   messageTypeGif,
}
