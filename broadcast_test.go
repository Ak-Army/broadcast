package broadcast

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BroadcastTestSuite struct {
	suite.Suite
}

func TestBroadcast(t *testing.T) {
	suite.Run(t, new(BroadcastTestSuite))
}

func (suite *BroadcastTestSuite) TestBroadcast() {
	received := make(chan int)
	s := New()
	receive := func(r Receiver, c chan int) {
		for {
			if m, ok := r.Receive(); ok {
				c <- m.(int)
				continue
			}
			return
		}
	}
	go receive(s.Receiver(), received)
	go receive(s.Receiver(), received)
	s.Send(2)
	suite.Equal(2, <-received)
	suite.Equal(2, <-received)
	s.Close()
}

func (suite *BroadcastTestSuite) TestSharedBroadcast() {
	received := make(chan int)
	f := func(r Receiver, c chan int) {
		for {
			if m, ok := r.Receive(); ok {
				c <- m.(int)
				continue
			}
			return
		}
	}
	ss := NewShared()
	go f(ss.Receiver(), received)
	go f(ss.Receiver(), received)
	go func() {
		ss.Send(2)
		ss.Send(1)
	}()
	go func() {
		ss.Send(1)
		ss.Send(2)
	}()
	ints := make([]int, 8)
	for i := range ints {
		n := <-received
		switch n {
		case 1, 2:
		default:
			suite.FailNow("want 1 or 2 got %v", n)
		}
		ints[i] = n
	}
	ss.Close()
}

func (suite *BroadcastTestSuite) TestBroadcastChanReceiver() {
	received := make(chan int)
	s := New()
	receive := func(r Receiver, c chan int) {
		for {
			select {
			case <-r.Chan():
				if m, ok := r.Val(); ok {
					c <- m.(int)
					continue
				}
			}
			return
		}
	}
	go receive(s.Receiver(), received)
	go receive(s.Receiver(), received)
	s.Send(2)
	suite.Equal(2, <-received)
	suite.Equal(2, <-received)
	s.Close()
}

/*

Run benchmarking with: go test -bench ./...

*/
func BenchmarkNSubscribers1MessageSender(b *testing.B) {
	wg := new(sync.WaitGroup)
	s := New()
	n := b.N
	if n > 8000 {
		n = 8000
	}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(r Receiver) {
			r.Receive()
			wg.Done()
		}(s.Receiver())
	}
	b.ResetTimer()
	s.Send(1)
	s.Close()
	wg.Wait()
}

func Benchmark1SubscriberNMessagesSender(b *testing.B) {
	wg := new(sync.WaitGroup)
	s := New()
	wg.Add(1)
	go func(r Receiver) {
		for i := 0; i < b.N; i++ {
			r.Receive()
		}
		wg.Done()
	}(s.Receiver())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Send(i)
	}
	s.Close()
	wg.Wait()
}

func BenchmarkNSubscriberNMessagesSender(b *testing.B) {
	wg := new(sync.WaitGroup)
	s := New()
	n := b.N
	if n > 8000 {
		n = 8000
	}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(r Receiver) {
			for {
				if _, ok := r.Receive(); !ok {
					break
				}
			}
			wg.Done()
		}(s.Receiver())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Send(i)
	}
	s.Close()
	wg.Wait()
}
