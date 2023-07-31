package stream

import (
	"fmt"
	"strconv"
	"testing"
)

func TestNewSliceSource(t *testing.T) {
	sink := &ChanSink{in: make(chan any)}
	NewSliceSource[int]([]int{1, 2, 3, 4, 5}).Link(NewFilterStream[int](func(i int) bool {
		if i%2 == 0 {
			return true
		}
		return false
	})).Link(NewMapStream[int, string](func(r int) string {
		return strconv.Itoa(r) + "_Map"
	})).Link(NewReduceStream[string](func(s string, s2 string) string {
		return s + s2
	})).To(sink)
	slice := SinkSlice[string](sink)
	fmt.Print(slice)
}

func TestMergeFlow(t *testing.T) {
	sink := NewChanSink()
	f1 := NewSliceSource[int]([]int{1, 2, 3, 4, 5})
	f2 := NewSliceSource[string]([]string{"a", "b", "c"})
	MergeFlow([]OutLet{f1, f2}).To(sink)
	s := SinkSlice[any](sink)
	fmt.Println(s)
}

func TestNewFileSource(t *testing.T) {
	sink := NewChanSink()
	f1 := NewSliceSource[int]([]int{1, 2, 3, 4, 5})
	f2, err := NewFileSource("./README.md")
	if err != nil {
		t.Log(err)
		return
	}
	f3 := f2.Link(NewMapStream[byte, string](func(b byte) string {
		return string(b)
	}))
	MergeFlow([]OutLet{f1, f3}).To(sink)
	s := SinkSlice[any](sink)
	fmt.Println(s)
}

func TestNewTcpSource(t *testing.T) {
	sink := NewChanSink()
	f, err := NewTcpSource("localhost:8888")
	if err != nil {
		t.Log(err)
		return
	}
	f.Link(NewBaseStream()).To(sink)
	_ = SinkSlice[byte](sink)
	//fmt.Println(string(s))
}

func TestNewStdinSource(t *testing.T) {
	// test中stdin=/dev/null
	sink := NewChanSink()
	NewStdinSource().Link(NewBaseStream()).To(sink)
	fmt.Println(SinkSlice[byte](sink))
}
