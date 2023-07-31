package stream

type InLet interface {
	In() chan any
}

type Sink interface {
	InLet
	Count() int
}

type ChanSink struct {
	in chan any
}

func (cs *ChanSink) In() chan any {
	return cs.in
}

func (cs *ChanSink) Count() int {
	cnt := 0
	for {
		_, ok := <-cs.In()
		if !ok {
			break
		}
		cnt++
	}
	return cnt
}

func NewChanSink() *ChanSink {
	return &ChanSink{in: make(chan any)}
}

func SinkSlice[T any](s Sink) []T {
	// sink node => slice
	res := make([]T, 0)
	for {
		v, ok := <-s.In()
		if !ok {
			break
		}
		res = append(res, v.(T))
	}
	return res
}

func SinkMap[T comparable, R any](s Sink) map[T]R {
	res := make(map[T]R)
	for {
		v, ok := <-s.In()
		if !ok {
			break
		}
		vp := v.(Pair[T, R])
		res[vp.First] = vp.Second
	}
	return res
}
