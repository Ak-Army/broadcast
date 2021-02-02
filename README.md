# Broadcast
A lock less non-blocking, concurrent safe broadcast application


# Example

For simple sender:
```go
package main

import (
	"fmt"

	"github.com/Ak-Army/broadcast"
)

func main() {
    received := make(chan int)
    s := broadcast.New()
    receive := func(r broadcast.Receiver, c chan int) {
        for {
            if m, ok := r.Receive(); ok {
                fmt.Println(m)
                continue
            }
            return
        }
    }
    go receive(s.Receiver(), received)
    go receive(s.Receiver(), received)
    s.Send(2)
    s.Send(11)
	s.Close()
}
```


For multi sender:
```go
package main

import (
	"fmt"

	"github.com/Ak-Army/broadcast"
)

func main() {
    received := make(chan int)
    ss := broadcast.NewShared()
    receive := func(r broadcast.Receiver, c chan int) {
        for {
            if m, ok := r.Receive(); ok {
                fmt.Println(m)
                continue
            }
            return
        }
    }
    go receive(ss.Receiver(), received)
    go receive(ss.Receiver(), received)
	go func() {
		ss.Send(22)
		ss.Send(12)
	}()
	go func() {
		ss.Send(11)
		ss.Send(21)
	}()
	ss.Close()
}
```
