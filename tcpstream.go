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
			continue
		}

		go t.do(srcConn, dstConn)
	}
}

func (t *Thread) do(src, dst io.ReadWriteCloser) {

	defer src.Close()
	defer dst.Close()
	go flow(src, dst, t.SrcToDstHook)
	flow(dst, src, t.DstToSrcHook)
}

func flow(src, dst io.ReadWriter, hook func(b []byte)) {
	buff := make([]byte, 0xffff)
	for {
		n, err := src.Read(buff)
		if err != nil {
			return
		}

		b := buff[:n]
		if hook != nil {
			hook(b)
		}

		dst.Write(b)
	}
}
