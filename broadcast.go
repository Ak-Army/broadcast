package broadcast

func New() Sender {
	return &sender{
		b: &broadcast{
			ch: make(chan struct{}),
		},
	}
}

func NewShared() Sender {
	return &sharedSender{
		sender: &sender{
			b: &broadcast{
				ch: make(chan struct{}),
			},
		},
	}
}
