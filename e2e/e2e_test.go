package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	rw := connect(t)
	write(t, rw, "PING")
	assert.Equal(t, "+PONG\r\n", read(t, rw))
}

func TestEcho(t *testing.T) {
	rw := connect(t)
	write(t, rw, "ECHO", "Hello")
	assert.Equal(t, "+Hello\r\n", read(t, rw))
}

func TestSetGet(t *testing.T) {
	rw := connect(t)
	write(t, rw, "SET", "Lewis", "Hamilton")
	assert.Equal(t, "+OK\r\n", read(t, rw))
	write(t, rw, "GET", "Lewis")
	assert.Equal(t, "+Hamilton\r\n", read(t, rw))
}

func read(t *testing.T, r *bufio.ReadWriter) string {
	s, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("Exiting due to error", err)
		t.Fatalf("Read error %v", err)
	}
	return s
}

func write(t *testing.T, w *bufio.ReadWriter, s ...string) {
	_, err := w.WriteString("*" + str(len(s)) + "\r\n")
	if err != nil {
		t.Errorf("write error %v", err)
	}
	for _, v := range s {
		_, err = w.WriteString("$" + str(len(v)) + "\r\n" + v + "\r\n")
		if err != nil {
			t.Errorf("write error %v", err)
		}
	}
	w.Flush()
}

func str(n int) string {
	return strconv.Itoa(n)
}

func connect(t *testing.T) *bufio.ReadWriter {
	conn, err := net.Dial("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Exiting due to error", err)
		t.Fatalf("Connect error %v", err)
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
}
