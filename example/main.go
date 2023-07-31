package main

import (
	"fmt"
	"github.com/cmz2012/stream"
)

func main() {
	sink := stream.NewChanSink()
	stream.NewStdinSource().Link(stream.NewBaseStream()).To(sink)
	fmt.Println(string(stream.SinkSlice[byte](sink)))
}
