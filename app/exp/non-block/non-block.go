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

	fmt.Println("Accepting connections")

	for {
		cfd, _, err := syscall.Accept(sfd)
		if err != nil {
			if isEagain(err) {
				continue
			}
			return err
		}
		err = syscall.SetNonblock(cfd, true)
		if err != nil {
			return err
		}
		go func() error {
			data := make([]byte, 2000)
			for {
				n, err := syscall.Read(cfd, data)
				if err != nil {
					if isEagain(err) {
						continue
					}
					return err
				}
				fmt.Printf("n : %d\n", n)
				if n > 0 {
					fmt.Print(string(data[:n]))
					n, err := syscall.Write(cfd, []byte("-ERR unknown command 'helloworld'\r\n"))
					fmt.Printf("Wrote %d bytes %v\n", n, err)
				} else {
					fmt.Println("Closing connection")
					return nil
				}
			}
		}()

	}
}

func isEagain(err error) bool {
	return err.(syscall.Errno).Timeout()
}
