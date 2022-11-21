package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type respType int

const (
	invalid respType = iota
	respString
	respError
	respInt
	respArrayLen
	respBulkStringLen
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

func parse(s string) (respType, interface{}, error) {
	s = strings.Trim(s, " ")

	if s[0] == ':' || s[0] == '*' || s[0] == '$' {
		n, err := strconv.Atoi(s[1:])
		if err != nil {
			return invalid, nil, fmt.Errorf("cannot parse string to int")
		}

		if s[0] == ':' {
			return respInt, n, nil
		} else if s[0] == '*' {
			return respArrayLen, n, nil
		} else {
			return respBulkStringLen, n, nil
		}
	}

	if s[0] == '+' {
		return respString, s[1:], nil
	}

	return invalid, nil, fmt.Errorf("unknown type")
}
