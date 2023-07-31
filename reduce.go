package stream

type ReduceStream[T any] struct {
	in  chan any
	out chan any
	f   Reduce[T]
}

func NewReduceStream[T any](f Reduce[T]) *ReduceStream[T] {
	s := &ReduceStream[T]{
		in:  make(chan any),
		out: make(chan any),
		f:   f,
	}
	s.Compute()
	return s
}

func (m *ReduceStream[T]) Link(f Flow) Flow {
	DoStream(m, f)
	return f
}

func (m *ReduceStream[T]) In() chan any {
	return m.in
}

func (m *ReduceStream[T]) Out() chan any {
	return m.out
}

func (m *ReduceStream[T]) Compute() {
	go func() {
		var rst T
		for {
			v1, ok := <-m.In()
			if !ok {
				break
			}
			rst = m.f(v1.(T), rst)
		}
		m.Out() <- rst
		close(m.Out())
	}()
}

func (m *ReduceStream[T]) To(s Sink) {
	DoStream(m, s)
}
