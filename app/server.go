package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	redis "redis-go/app/redis_go"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	chcmd := make(chan redis.Command)

	go start(l, chcmd)
	process(chcmd)
}

func process(chcmd chan redis.Command) {
	data := make(map[string]string)
	for {
		c := <-chcmd
		c.Response() <- c.Execute(&data)
	}
}

func start(l net.Listener, chcmd chan redis.Command) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go execute(conn, chcmd)
	}
}

func execute(conn net.Conn, chcmd chan redis.Command) {
	rr := redis.NewRespReader(bufio.NewReader(conn))
	cr := redis.NewCommandReader(rr)

	for {
		c, err := cr.Read()
		if err != nil {
			fmt.Printf("Error reading command: %v\n", err)
			conn.Close()
			return
		}
		chcmd <- c
		resp := <-c.Response()
		conn.Write([]byte(resp + "\r\n"))
	}
}
