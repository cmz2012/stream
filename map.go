package stream

type MapStream[T, R any] struct {
	in  chan any
	out chan any
	f   Map[T, R]
}

func NewMapStream[T, R any](f Map[T, R]) *MapStream[T, R] {
	s := &MapStream[T, R]{
		in:  make(chan any),
		out: make(chan any),
		f:   f,
	}
	s.Compute()
	return s
}

func (m *MapStream[T, R]) Link(f Flow) Flow {
	DoStream(m, f)
	return f
}

func (m *MapStream[T, R]) In() chan any {
	return m.in
}

func (m *MapStream[T, R]) Out() chan any {
	return m.out
}

func (m *MapStream[T, R]) Compute() {
	go func() {
		for {
			v, ok := <-m.In()
			if !ok {
				break
			}
			m.Out() <- m.f(v.(T))
		}
		close(m.Out())
	}()
}

func (m *MapStream[T, R]) To(s Sink) {
	DoStream(m, s)
}
