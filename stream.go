package stream

type Flow interface {
	Link(f Flow) Flow
	InLet
	OutLet
	To(s Sink)
}

func DoStream(from OutLet, in InLet) {
	go func() {
		for {
			v, ok := <-from.Out()
			if !ok {
				break
			}
			in.In() <- v
		}
		close(in.In())
	}()
}
