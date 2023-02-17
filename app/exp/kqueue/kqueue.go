package main

import (
	"fmt"
	"syscall"
)

func main() {
	err := connect()
	if err != nil {
		fmt.Println(err)
	}
}

func connect() error {

	sfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(sfd)

	var sa syscall.SockaddrInet4
	sa.Port = 6379
	sa.Addr = [4]byte{0, 0, 0, 0}

	err = syscall.Bind(sfd, &sa)
	if err != nil {
		return err
	}

	err = syscall.Listen(sfd, 3)
	if err != nil {
		return err
	}

	err = syscall.SetNonblock(sfd, true)
	if err != nil {
		return err
	}

	kq, err := syscall.Kqueue()
	if err != nil {
		return err
	}

	err = addEvent(kq, sfd)
	if err != nil {
		return err
	}

	fmt.Println("Accepting connections")

	for {
		events := make([]syscall.Kevent_t, 10)
		n, err := syscall.Kevent(kq, nil, events, nil)
		if err != nil {
			fmt.Printf("Error waiting for kqueue events: %v\n", err)
			continue
		}

		for i := 0; i < n; i++ {
			if events[i].Ident == uint64(sfd) {
				err = accept(&events[i], kq, sfd)
				if err != nil {
					return err
				}
			} else {
				cfd := int(events[i].Ident)
				err = read(cfd)
				if err != nil {
					return err
				}
			}
		}
	}
}

func read(cfd int) error {
	data := make([]byte, 2000)
	n, err := syscall.Read(cfd, data)
	if err != nil {
		return err
	}
	fmt.Printf("n : %d\n", n)
	if n > 0 {
		fmt.Print(string(data[:n]))
	}
	// if n == 0, we should remove from kq?
	return nil
}

func accept(ev *syscall.Kevent_t, kq int, sfd int) error {
	cfd, _, err := syscall.Accept(sfd)
	if err != nil {
		return err
	}
	err = syscall.SetNonblock(cfd, true)
	if err != nil {
		return err
	}

	err = addEvent(kq, cfd)
	if err != nil {
		return err
	}

	return nil
}

func addEvent(kq int, fd int) error {
	ev := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE | syscall.EV_CLEAR,
	}

	_, err := syscall.Kevent(kq, []syscall.Kevent_t{ev}, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
