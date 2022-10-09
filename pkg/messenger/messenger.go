package messenger

type Message struct {
	Conversation int64
	Text         string
}

type Messenger interface {
	Message(Message) error
}
