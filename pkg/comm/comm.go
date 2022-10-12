package comm

type Message struct {
	Conversation int64
	Text         string
}

type Messenger interface {
	Message(Message) error
}

type Handler func(Message) error

type Listener interface {
	Listen(Handler) error
}
