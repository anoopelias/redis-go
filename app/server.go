package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"redis-go/app/resp"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go execute(conn)
	}
}

func execute(conn net.Conn) {
	rr := resp.NewRespReader(bufio.NewReader(conn))
	cr := resp.NewCommandReader(rr)

	for {
		// TODO: Handle error
		cr.Read()
		conn.Write([]byte("+PONG\r\n"))
	}
}
