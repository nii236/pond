package pond

type Message struct {
	Message string
	Meta    map[string]string
}

type BotAdaptor interface {
	Init(chan *Message, chan *Message) error
	Run()
}
