package comm

type Message struct {
	Conversation int64
	Text         string
}

//go:generate mockgen -destination=mocks/messenger.go -package=mocks . Messenger
type Messenger interface {
	Message(Message) error
}

type Handler func(Message) error

type Listener interface {
	Listen(Handler) error
}
