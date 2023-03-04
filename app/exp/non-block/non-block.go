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

	fmt.Println("Accepting connections")
	err = loop(sfd)
	if err != nil {
		panic(err)
	}
}

func loop(sfd int) error {
	cfds := []int{}
	dict := make(map[string]string)

	for {
		// Handle new connections
		cfd, isNew, err := accept(sfd)
		if err != nil {
			return err
		}
		if isNew {
			cfds = append(cfds, cfd)
		}

		ncfds := []int{}
		// Check existing connections
		for _, cfd := range cfds {
			ctd, err := handle(cfd, dict)
			if err != nil {
				return err
			}
			if ctd {
				ncfds = append(ncfds, cfd)
			}
		}
		cfds = ncfds
	}
}

func accept(sfd int) (int, bool, error) {
	isNew := true
	cfd, _, err := syscall.Accept(sfd)
	if err != nil {
		if shouldRetry(err) {
			isNew = false
		} else {
			return -1, isNew, err
		}
	}

	if isNew {
		err = syscall.SetNonblock(cfd, true)
		if err != nil {
			return -1, isNew, err
		}
	}
	return cfd, isNew, nil
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

func handle(cfd int, dict map[string]string) (bool, error) {
	ctd := true
	data := make([]byte, 2000)
	n, err := syscall.Read(cfd, data)
	if err != nil {
		if shouldRetry(err) {
			return ctd, nil
		}
		return ctd, err
	}
	if n > 0 {
		req := string(data[:n])
		splits := strings.Split(req, "\r\n")

		if strings.EqualFold(splits[2], "GET") {
			err = get(cfd, splits[4], dict)
			if err != nil {
				return ctd, err
			}
		} else if strings.EqualFold(splits[2], "SET") {
			err = set(cfd, splits[4], splits[6], dict)
			if err != nil {
				return ctd, err
			}
		} else {
			// We just say OK for unknown commands
			_, err := syscall.Write(cfd, []byte("+OK\r\n"))
			if err != nil {
				return ctd, err
			}
		}
	} else if n == 0 {
		err = syscall.Close(cfd)
		if err != nil {
			return ctd, err
		}
		return false, nil
	}
	return ctd, nil
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

func shouldRetry(err error) bool {
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	return errno.Temporary()
}
