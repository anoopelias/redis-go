package main

import (
	"fmt"
	"strings"
	"syscall"
)

func main() {
	err := connect()
	if err != nil {
		fmt.Println(err)
	}
}

func connect() error {

	dict := make(map[string]string)

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

	err = syscall.Listen(sfd, 50)
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
		events := make([]syscall.Kevent_t, 100)
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
				err = handle(cfd, dict)
				if err != nil {
					return err
				}
			}
		}
	}
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
		if splits[2] == "GET" {
			err = get(cfd, splits[4], dict)
			if err != nil {
				return err
			}
		} else if splits[2] == "SET" {
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
