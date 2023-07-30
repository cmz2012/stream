package stream

type BaseStream struct {
	in  chan any
	out chan any
}

func NewBaseStream() *BaseStream {
	s := &BaseStream{
		in:  make(chan any),
		out: make(chan any),
	}
	s.Compute()
	return s
}

func (m *BaseStream) Link(f Flow) Flow {
	DoStream(m, f)
	return f
}

func (m *BaseStream) In() chan any {
	return m.in
}

func (m *BaseStream) Out() chan any {
	return m.out
}

func (m *BaseStream) Compute() {
	go func() {
		for {
			v, ok := <-m.In()
			if !ok {
				break
			}
			m.Out() <- v
		}
		close(m.Out())
	}()
}

func (m *BaseStream) To(s Sink) {
	DoStream(m, s)
}
