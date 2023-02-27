package ev

import (
	"fmt"
	"syscall"
)

type SysCall interface {
	Socket(int, int, int) (int, error)
	Bind(int, syscall.Sockaddr) error
	Listen(int, int) error
	SetNonblock(int, bool) error
	Accept(int) (int, syscall.Sockaddr, error)
	Read(int, []byte) (int, error)
	Write(int, []byte) (int, error)
	Close(int) error
}

type SysCallError interface {
	Error() string
	Timeout() bool
}

type EventLoop interface {
	Run(func([]byte) string) error
}

type SocketEventLoop struct {
	handler func(StringReader) string
	cfds    []int
	sys     SysCall
}

func NewSocketEventLoop(sys SysCall) SocketEventLoop {
	return SocketEventLoop{
		cfds: make([]int, 0),
		sys:  sys,
	}
}

func (el *SocketEventLoop) Run(handler func(StringReader) string) error {
	el.handler = handler
	sfd, err := el.create()
	if err != nil {
		return err
	}
	fmt.Println("Server started")

	for {
		err = el.execute(sfd)
		if err != nil {
			return err
		}
	}
}

func (el *SocketEventLoop) create() (int, error) {
	sfd, err := el.sys.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return -1, err
	}

	// TODO: Port needs to be configurable
	sa := syscall.SockaddrInet4{
		Port: 6379,
		Addr: [4]byte{0, 0, 0, 0},
	}

	err = el.sys.Bind(sfd, &sa)
	if err != nil {
		return -1, err
	}

	// Setting backlog to something acceptable to redis-benchmark command
	// Needs more thoughts on what should be the ideal value here.
	err = el.sys.Listen(sfd, 50)
	if err != nil {
		return -1, err
	}

	err = el.sys.SetNonblock(sfd, true)
	if err != nil {
		return -1, err
	}
	return sfd, nil
}

func (el *SocketEventLoop) execute(sfd int) error {
	err := el.accept(sfd)
	if err != nil {
		return err
	}

	cfds := []int{}
	for _, cfd := range el.cfds {
		ctd, err := el.process(cfd)
		if err != nil {
			return err
		}

		if ctd {
			cfds = append(cfds, cfd)
		} else {
			err = el.sys.Close(cfd)
			if err != nil {
				return err
			}
		}

	}
	el.cfds = cfds

	return nil
}

func (el *SocketEventLoop) accept(sfd int) error {
	cfd, _, err := el.sys.Accept(sfd)
	isNew := true
	if err != nil {
		isNew = false
		if !shouldRetry(err) {
			return err
		}
	}

	if isNew {
		err = el.sys.SetNonblock(cfd, true)
		if err != nil {
			return err
		}
		el.cfds = append(el.cfds, cfd)
	}
	return nil
}

func (el *SocketEventLoop) process(cfd int) (bool, error) {
	ctd := true
	data, err := el.read(cfd)
	if err != nil {
		if shouldRetry(err) {
			return ctd, nil
		}
		return ctd, err
	}

	if len(data) == 0 {
		return false, nil
	}

	sr := NewArrayStringReader(data)
	res := el.handler(sr)

	data = []byte(res)
	n, err := el.sys.Write(cfd, data)
	if err != nil {
		return ctd, err
	}

	if n != len(data) {
		// TODO: Retry rest of the bytes
		return ctd, fmt.Errorf("incomplete write")
	}

	return ctd, nil
}

func (el *SocketEventLoop) read(cfd int) ([]byte, error) {
	// TODO: Find the correct size for this buffer
	data := make([]byte, 2000)
	n, err := el.sys.Read(cfd, data)
	if err != nil {
		return data, err
	}
	return data[:n], nil
}

func shouldRetry(err error) bool {
	sce, ok := err.(SysCallError)
	if !ok {
		return false
	}
	return sce.Timeout()
}
