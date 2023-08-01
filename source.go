package stream

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type OutLet interface {
	Out() chan any
}

type Source interface {
	OutLet
	Link(f Flow) Flow
}

type ChanSource struct {
	out chan any
}

func (cs *ChanSource) Out() chan any {
	return cs.out
}

func (cs *ChanSource) Link(f Flow) Flow {
	DoStream(cs, f)
	return f
}

// NewSliceSource 切片source
func NewSliceSource[T any](slice []T) Source {
	// slice => source node
	out := make(chan any)
	source := &ChanSource{out}
	go func() {
		for _, v := range slice {
			source.Out() <- v
		}
		close(source.Out())
	}()
	return source
}

func NewMapSource[T comparable, R any](m map[T]R) Source {
	source := &ChanSource{out: make(chan any)}
	go func() {
		for k, v := range m {
			source.Out() <- Pair[T, R]{
				First:  k,
				Second: v,
			}
		}
		close(source.Out())
	}()
	return source
}

type FileSource struct {
	f   *os.File
	out chan any
}

func (fs *FileSource) Out() chan any {
	return fs.out
}

func (fs *FileSource) Link(f Flow) Flow {
	DoStream(fs, f)
	return f
}

// NewFileSource 文件source
func NewFileSource(f *os.File) (s Source) {
	s = &FileSource{
		f:   f,
		out: make(chan any),
	}
	go func() {
		defer f.Close()
		reader := bufio.NewReader(f)
		for {
			b, err := reader.ReadByte()
			if err != nil {
				fmt.Printf("File reader: %v\n", err)
				break
			}
			s.Out() <- b
		}
		close(s.Out())
	}()
	return
}

type TcpSource struct {
	listener *net.TCPListener
	out      chan any
}

func (ts *TcpSource) Out() chan any {
	return ts.out
}

func (ts *TcpSource) Link(f Flow) Flow {
	DoStream(ts, f)
	return f
}

// NewTcpSource TCP网络source
func NewTcpSource(laddr string) (s Source, err error) {
	addr, _ := net.ResolveTCPAddr("tcp", laddr)
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	s = &TcpSource{
		listener: listener,
		out:      make(chan any),
	}
	go func() {
		for {
			log.Println("start listen...")
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("AcceptTCP: %v\n", err)
				break
			}
			log.Println("accept conn:  ", conn.RemoteAddr().String())
			r := bufio.NewReader(conn)
			//w := bufio.NewWriter(conn)
			for {
				b, err := r.ReadByte()
				if err != nil {
					log.Printf("Read Done: %v, %v\n", err, conn.RemoteAddr().String())
					break
				}
				s.Out() <- b
			}
			conn.Close()
			close(s.Out())
		}
	}()
	return
}
