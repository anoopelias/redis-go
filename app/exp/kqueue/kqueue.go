package main

import (
	"fmt"
	"strings"
	"syscall"
)

func main() {

	sfd, err := startServer()
	if err != nil {
		panic(err)
	}
	defer syscall.Close(sfd)

	// Create kqueue
	kq, err := syscall.Kqueue()
	if err != nil {
		panic(err)
	}

	// Add sfd to kqueue
	err = addEvent(kq, sfd)
	if err != nil {
		panic(err)
	}

	fmt.Println("Accepting connections")

	err = loop(kq, sfd)
	if err != nil {
		panic(err)
	}
}

func loop(kq int, sfd int) error {
	dict := make(map[string]string)
	for {
		events := make([]syscall.Kevent_t, 100)
		n, err := syscall.Kevent(kq, nil, events, nil)
		// There is a possible EINTR for which we need to retry.
		if err != nil && !shouldRetry(err) {
			return err
		}

		for i := 0; i < n; i++ {
			if events[i].Ident == uint64(sfd) {
				err = accept(&events[i], kq, sfd)
				if err != nil {
					return err
				}
			} else {
				cfd := int(events[i].Ident)
				err = handle(cfd, dict)
				if err != nil {
					return err
				}
			}
		}
	}
}

func startServer() (int, error) {
	sfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return 0, err
	}

	var sa syscall.SockaddrInet4
	sa.Port = 6379
	sa.Addr = [4]byte{0, 0, 0, 0}

	err = syscall.Bind(sfd, &sa)
	if err != nil {
		return 0, err
	}

	err = syscall.Listen(sfd, 50)
	if err != nil {
		return 0, err
	}

	err = syscall.SetNonblock(sfd, true)
	if err != nil {
		return 0, err
	}
	return sfd, nil
}

func handle(cfd int, dict map[string]string) error {
	data := make([]byte, 2000)
	n, err := syscall.Read(cfd, data)
	if err != nil {
		return err
	}
	if n > 0 {
		req := string(data[:n])
		splits := strings.Split(req, "\r\n")

		if strings.EqualFold(splits[2], "GET") {
			err = get(cfd, splits[4], dict)
			if err != nil {
				return err
			}
		} else if strings.EqualFold(splits[2], "SET") {
			err = set(cfd, splits[4], splits[6], dict)
			if err != nil {
				return err
			}
		} else {
			// We just say OK for unknown commands
			_, err := syscall.Write(cfd, []byte("+OK\r\n"))
			if err != nil {
				return err
			}
		}
	} else if n == 0 {
		err = syscall.Close(cfd)
		if err != nil {
			return err
		}
	}
	return nil
}

func get(cfd int, key string, dict map[string]string) error {
	if v, ok := dict[key]; ok {
		_, err := syscall.Write(cfd, []byte("+"+v+"\r\n"))
		if err != nil {
			return err
		}
	} else {
		_, err := syscall.Write(cfd, []byte("$-1\r\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

func set(cfd int, key string, value string, dict map[string]string) error {
	dict[key] = value
	_, err := syscall.Write(cfd, []byte("+OK\r\n"))
	if err != nil {
		return err
	}

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
		Ident: uint64(fd),
		// Filter read operations
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD,
	}

	_, err := syscall.Kevent(kq, []syscall.Kevent_t{ev}, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func shouldRetry(err error) bool {
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	return errno.Temporary()
}
