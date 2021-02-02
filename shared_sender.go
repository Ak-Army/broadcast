package broadcast

import "sync"

type sharedSender struct {
	*sender
	m sync.Mutex
}

func (s *sharedSender) Send(b Message) {
	s.m.Lock()
	defer s.m.Unlock()

	s.sender.Send(b)
}

func (s *sharedSender) Receiver() Receiver {
	s.m.Lock()
	defer s.m.Unlock()

	return s.sender.Receiver()
}

func (s *sharedSender) Close() {
	s.m.Lock()
	defer s.m.Unlock()

	s.sender.Close()
}
