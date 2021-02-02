package broadcast

type sender struct {
	b *broadcast
}

type broadcast struct {
	ch   chan struct{}
	val  Message
	next *broadcast
}

func (s *sender) Send(b Message) {
	s.b = s.b.send(b)
}

func (s *sender) Receiver() Receiver {
	return &receiver{b: s.b}
}

func (s *sender) Close() {
	s.b.close()
}

func (b *broadcast) send(bc Message) *broadcast {
	b.val = bc
	n := &broadcast{
		ch: make(chan struct{}),
	}
	b.next = n
	close(b.ch)

	return n
}

func (b *broadcast) close() {
	b.val = nil
	b.next = nil
	close(b.ch)
}
