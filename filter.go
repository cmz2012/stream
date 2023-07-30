package stream

type FilterStream[T any] struct {
	in  chan any
	out chan any
	f   Filter[T]
}

func NewFilterStream[T any](f Filter[T]) *FilterStream[T] {
	s := &FilterStream[T]{
		in:  make(chan any),
		out: make(chan any),
		f:   f,
	}
	s.Compute()
	return s
}

func (m *FilterStream[T]) Link(f Flow) Flow {
	DoStream(m, f)
	return f
}

func (m *FilterStream[T]) In() chan any {
	return m.in
}

func (m *FilterStream[T]) Out() chan any {
	return m.out
}

func (m *FilterStream[T]) Compute() {
	go func() {
		for {
			v, ok := <-m.In()
			if !ok {
				break
			}
			if m.f(v.(T)) {
				m.Out() <- v.(T)
			}
		}
		close(m.Out())
	}()
}

func (m *FilterStream[T]) To(s Sink) {
	DoStream(m, s)
}
