package stream

import (
	"code.byted.org/lang/gg/gslice"
)

type SortStream[T any] struct {
	in     chan any
	out    chan any
	f      SortBy[T]
	offset int
	limit  int
}

func NewSortStream[T any](f SortBy[T], offset, limit int) *SortStream[T] {
	s := &SortStream[T]{
		in:     make(chan any),
		out:    make(chan any),
		f:      f,
		offset: offset,
		limit:  limit,
	}
	s.Compute()
	return s
}

func (m *SortStream[T]) Link(f Flow) Flow {
	DoStream(m, f)
	return f
}

func (m *SortStream[T]) In() chan any {
	return m.in
}

func (m *SortStream[T]) Out() chan any {
	return m.out
}

func (m *SortStream[T]) Compute() {
	go func() {
		data := make([]T, 0)
		for {
			v1, ok := <-m.In()
			if !ok {
				break
			}
			data = append(data, v1.(T))
		}
		gslice.SortBy(data, m.f)
		data = data[m.offset : m.offset+m.limit]
		for _, d := range data {
			m.Out() <- d
		}
		close(m.Out())
	}()
}

func (m *SortStream[T]) To(s Sink) {
	DoStream(m, s)
}
