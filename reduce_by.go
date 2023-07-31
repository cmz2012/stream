package stream

type ReduceByStream[T comparable, R any] struct {
	in  chan any
	out chan any
	f   Reduce[R]
}

func NewReduceByStream[T comparable, R any](f Reduce[R]) *ReduceByStream[T, R] {
	s := &ReduceByStream[T, R]{
		in:  make(chan any),
		out: make(chan any),
		f:   f,
	}
	s.Compute()
	return s
}

func (m *ReduceByStream[T, R]) Link(f Flow) Flow {
	DoStream(m, f)
	return f
}

func (m *ReduceByStream[T, R]) In() chan any {
	return m.in
}

func (m *ReduceByStream[T, R]) Out() chan any {
	return m.out
}

func (m *ReduceByStream[T, R]) Compute() {
	go func() {
		reduceMap := make(map[T]R)
		for {
			v1, ok := <-m.In()
			if !ok {
				break
			}
			vp := v1.(Pair[T, R])
			first, second := vp.First, vp.Second
			reduceMap[first] = m.f(reduceMap[first], second)
		}
		for k, v := range reduceMap {
			m.Out() <- Pair[T, R]{First: k, Second: v}
		}
		close(m.Out())
	}()
}

func (m *ReduceByStream[T, R]) To(s Sink) {
	DoStream(m, s)
}
