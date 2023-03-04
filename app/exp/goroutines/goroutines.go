package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type command struct {
	typ   int
	key   string
	value string
	resp  chan string
}

const (
	unknown int = iota
	set
	get
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		panic(err)
	}
	fmt.Println("Accepting connections")
	chcmd := make(chan command)

	go start(l, chcmd)
	loop(chcmd)
}

func loop(chcmd chan command) {
	dict := make(map[string]string)
	for {
		c := <-chcmd
		c.resp <- execute(c, dict)
	}
}

func execute(cmd command, dict map[string]string) string {
	switch cmd.typ {
	case get:
		return dict[cmd.key]
	case set:
		dict[cmd.key] = cmd.value
	}

	return "OK"
}

func start(l net.Listener, chcmd chan command) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			panic(err)
		}
		go handle(conn, chcmd)
	}
}

func handle(conn net.Conn, chcmd chan command) {
	r := bufio.NewReader(conn)
	for {
		c, err := readCommand(r)
		if err != nil {
			// Possibly because the client disconnected
			return
		}
		chcmd <- c
		resp := <-c.resp
		conn.Write([]byte("+" + resp + "\r\n"))
	}
}

func readCommand(r *bufio.Reader) (command, error) {
	cmd := command{}
	cmd.resp = make(chan string)

	// Length
	by, _, err := r.ReadLine()
	if err != nil {
		return cmd, err
	}
	n, _ := strconv.Atoi(string(by[1:]))

	for i := 0; i < n; i++ {
		_, _, err = r.ReadLine()
		if err != nil {
			return cmd, err
		}

		by, _, err = r.ReadLine()
		if err != nil {
			return cmd, err
		}

		l := string(by)
		switch i {
		case 0:
			if strings.EqualFold(l, "SET") {
				cmd.typ = set
			} else if strings.EqualFold(l, "GET") {
				cmd.typ = get
			}
		case 1:
			cmd.key = l
		case 2:
			cmd.value = l
		}
	}

	return cmd, nil

}
