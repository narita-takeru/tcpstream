package tcpstream

import (
	"fmt"
	"io"
	"net"
)

type Thread struct {
	SrcToDstHook func(b []byte)
	DstToSrcHook func(b []byte)
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

	for {
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

		go t.do(srcConn, dstConn)
	}
}

func (t *Thread) do(src, dst io.ReadWriteCloser) {

	defer src.Close()
	defer dst.Close()

	done := make(chan struct{}, 0)

	wk := worker{}
	go wk.flow(src, dst, t.SrcToDstHook, done)
	go wk.flow(dst, src, t.DstToSrcHook, done)

	<-done
}

type worker struct {
	closed bool
}

func (wk *worker) flow(src, dst io.ReadWriter, hook func(b []byte), done chan struct{}) {

	buff := make([]byte, 0xffff)
	for {
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
			hook(b)
		}

		dst.Write(b)
	}
}
