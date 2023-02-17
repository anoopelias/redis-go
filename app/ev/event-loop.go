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
	Kqueue() (int, error)
	Kevent(int, []syscall.Kevent_t, []syscall.Kevent_t, *syscall.Timespec) (n int, err error)
}

type SysCallError interface {
	Error() string
	Temporary() bool
}

type EventLoop interface {
	Run(func([]byte) string) error
}

type SocketEventLoop struct {
	handler func(StringReader) string
	cfds    []int
	sys     SysCall
	kq      int
	sfd     int
}

func NewSocketEventLoop(sys SysCall) SocketEventLoop {
	return SocketEventLoop{
		sys: sys,
	}
}

func (el *SocketEventLoop) Run(handler func(StringReader) string) error {
	el.handler = handler
	err := el.create()
	if err != nil {
		return err
	}
	fmt.Println("Server started")

	for {
		err = el.execute()
		if err != nil {
			return err
		}
	}
}

func (el *SocketEventLoop) addKqEvent(fd int) error {
	ev := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD,
	}
	_, err := el.sys.Kevent(el.kq, []syscall.Kevent_t{ev}, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (el *SocketEventLoop) create() error {
	sfd, err := el.sys.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	el.sfd = sfd

	// TODO: Port needs to be configurable
	sa := syscall.SockaddrInet4{
		Port: 6379,
		Addr: [4]byte{0, 0, 0, 0},
	}

	err = el.sys.Bind(sfd, &sa)
	if err != nil {
		return err
	}

	// Setting backlog to something acceptable to redis-benchmark command
	// Needs more thoughts on what should be the ideal value here.
	err = el.sys.Listen(sfd, 50)
	if err != nil {
		return err
	}

	err = el.sys.SetNonblock(sfd, true)
	if err != nil {
		return err
	}

	kq, err := el.sys.Kqueue()
	if err != nil {
		return err
	}
	el.kq = kq

	err = el.addKqEvent(sfd)
	if err != nil {
		return err
	}

	return nil
}

func (el *SocketEventLoop) execute() error {
	events := make([]syscall.Kevent_t, 10)
	n, err := el.sys.Kevent(el.kq, nil, events, nil)

	if err != nil && !shouldRetry(err) {
		return err
	}

	for i := 0; i < n; i++ {
		fid := int(events[i].Ident)

		if fid == el.sfd {
			err = el.accept()
			if err != nil {
				return err
			}
		} else {
			ctd, _ := el.process(fid)

			if !ctd {
				err = el.sys.Close(fid)
				if err != nil {
					return err
				}
			}

		}
	}
	return nil
}

func (el *SocketEventLoop) accept() error {
	cfd, _, err := el.sys.Accept(el.sfd)
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
		err := el.addKqEvent(cfd)
		if err != nil {
			return err
		}
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
	return sce.Temporary()
}
