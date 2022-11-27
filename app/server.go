package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"redis-go/app/redis_go"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	data := make(map[string]string)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go execute(conn, &data)
	}
}

func execute(conn net.Conn, data *map[string]string) {
	rr := redis_go.NewRespReader(bufio.NewReader(conn))
	cr := redis_go.NewCommandReader(rr)

	for {
		c, err := cr.Read()
		if err != nil {
			fmt.Printf("Error reading command: %v\n", err)
			conn.Close()
			return
		}
		conn.Write([]byte(c.Execute(data) + "\r\n"))
	}
}
