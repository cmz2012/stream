package main

import (
	"fmt"
	"github.com/cmz2012/stream"
	"os"
)

func main() {
	sink := stream.NewChanSink()
	stream.NewFileSource(os.Stdin).Link(stream.NewBaseStream()).To(sink)
	fmt.Println(string(stream.SinkSlice[byte](sink)))
}
