package stream

import (
	"reflect"
)

func SplitFlow[T any](f Flow, split Split[T]) []Flow {
	f1, f2 := NewBaseStream(), NewBaseStream()
	go func() {
		for {
			v, ok := <-f.Out()
			if !ok {
				break
			}
			if split(v.(T)) {
				f1.In() <- v
			} else {
				f2.In() <- v
			}
		}
		close(f1.In())
		close(f2.In())
	}()
	return []Flow{f1, f2}
}

func MergeFlow(fs []OutLet) Flow {
	nf := NewBaseStream()
	go func() {
		cases := make([]reflect.SelectCase, 0)
		for _, f := range fs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(f.Out()),
			})
		}
		cnt := len(cases)
		flg := make(map[int]bool)
		for len(flg) < cnt {
			//fmt.Printf("Select case: %v\n", flg)
			index, v, ok := reflect.Select(cases)
			if !ok {
				flg[index] = true
				continue
			}
			nf.In() <- v.Interface()
		}
		close(nf.In())
	}()
	return nf
}
