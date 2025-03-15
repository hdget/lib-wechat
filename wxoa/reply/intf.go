package reply

type ReceivedManager interface {
	Handle() ([]byte, error)
	ReplyText(content string) ([]byte, error)
}
