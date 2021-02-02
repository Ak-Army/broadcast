package broadcast

type Message interface{}

type Sender interface {
	Send(Message)
	Receiver() Receiver
	Close()
}

type Receiver interface {
	Receive() (Message, bool)
	Chan() <-chan struct{}
	Val() (Message, bool)
}
