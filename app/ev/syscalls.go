package ev

import "syscall"

type Syscalls struct{}

func (*Syscalls) Socket(domain, typ, proto int) (int, error) {
	return syscall.Socket(domain, typ, proto)
}

func (*Syscalls) Bind(fd int, sa syscall.Sockaddr) error {
	return syscall.Bind(fd, sa)
}

func (*Syscalls) Listen(s int, backlog int) error {
	return syscall.Listen(s, backlog)
}

func (*Syscalls) SetNonblock(fd int, nonblocking bool) error {
	return syscall.SetNonblock(fd, nonblocking)
}

func (*Syscalls) Accept(fd int) (int, syscall.Sockaddr, error) {
	return syscall.Accept(fd)
}

func (*Syscalls) Read(fd int, p []byte) (int, error) {
	return syscall.Read(fd, p)
}

func (*Syscalls) Write(fd int, p []byte) (int, error) {
	return syscall.Write(fd, p)
}

func (*Syscalls) Close(fd int) error {
	return syscall.Close(fd)
}

func (*Syscalls) Kqueue() (int, error) {
	return syscall.Kqueue()
}

func (*Syscalls) Kevent(kq int, changes, events []syscall.Kevent_t, timeout *syscall.Timespec) (n int, err error) {
	return syscall.Kevent(kq, changes, events, timeout)
}
