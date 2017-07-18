package tcpstream

import (
	"fmt"
	"io"
	"net"
)

type Thread struct {
	SrcToDstHook func(id, seq int, b []byte) []byte
	DstToSrcHook func(id, seq int, b []byte) []byte
}

func (t *Thread) Do(src, dst string) {

	srcAddr, err := net.ResolveTCPAddr("tcp", src)
	if err != nil {
		panic(err)
	}

	dstAddr, err := net.ResolveTCPAddr("tcp", dst)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", srcAddr)
	if err != nil {
		panic(err)
	}

	for i := 0; true; i++ {
		srcConn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		dstConn, err := net.DialTCP("tcp", nil, dstAddr)
		if err != nil {
			srcConn.Close()
			continue
		}

		go t.do(i, srcConn, dstConn)
	}
}

func (t *Thread) do(id int, src, dst io.ReadWriteCloser) {

	defer src.Close()
	defer dst.Close()

	done := make(chan struct{}, 0)

	wk := worker{}
	go wk.flow(id, src, dst, t.SrcToDstHook, done)
	go wk.flow(id, dst, src, t.DstToSrcHook, done)

	<-done
}

type worker struct {
	closed bool
}

func (wk *worker) flow(id int, src, dst io.ReadWriter, hook func(id, seq int, b []byte) []byte, done chan struct{}) {

	buff := make([]byte, 0xffff)
	for i := 0; true; i++ {
		n, err := src.Read(buff)
		if err != nil {
			if !wk.closed {
				wk.closed = true
				close(done)
			}

			return
		}

		b := buff[:n]
		if hook != nil {
			b = hook(id, i, b)
		}

		dst.Write(b)
	}
}
