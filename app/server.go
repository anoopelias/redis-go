package main

import (
	"bufio"
	"fmt"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	r := bufio.NewReader(conn)

	for {
		nstr := readStringLine(r)
		fmt.Println("NumberOfCommands: ", nstr)
		cnstr := readStringLine(r)
		fmt.Println("NumberOfChars: ", cnstr)
		c := readStringLine(r)
		fmt.Println("Command: ", c)

		conn.Write([]byte("+PONG\r\n"))
	}
}

func readStringLine(reader *bufio.Reader) string {
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading buffer: ", err.Error())
		os.Exit(1)
	}
	return line
}
