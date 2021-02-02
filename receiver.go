package broadcast

type receiver struct {
	b *broadcast
}

func (r *receiver) Receive() (Message, bool) {
	<-r.Chan()
	return r.Val()
}

func (r *receiver) Chan() <-chan struct{} {
	return r.b.ch
}

func (r *receiver) Val() (Message, bool) {
	val := r.b.val
	r.b = r.b.next
	return val, r.b != nil
}
