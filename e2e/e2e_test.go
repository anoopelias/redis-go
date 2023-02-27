package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	rw, err := connect()
	if err != nil {
		t.Error(err)
	}
	write(t, rw, "PING")
	assert.Equal(t, "PONG", read(t, rw))
}

func TestEcho(t *testing.T) {
	rw, err := connect()
	if err != nil {
		t.Error(err)
	}
	write(t, rw, "ECHO", "Hello")
	assert.Equal(t, "Hello", read(t, rw))
}

func TestSetGet(t *testing.T) {
	rw, err := connect()
	if err != nil {
		t.Error(err)
	}
	write(t, rw, "SET", "Lewis", "Hamilton")
	assert.Equal(t, "OK", read(t, rw))
	write(t, rw, "GET", "Lewis")
	assert.Equal(t, "Hamilton", read(t, rw))
}

type resp struct {
	connectTime int
	getSetTime  int
	err         error
}

func TestSetGetMulti(t *testing.T) {
	rchan := make(chan resp)

	n := 500
	for i := 0; i < n; i++ {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		go testSetGetKey(t, rchan, i)
	}

	ctt := 0
	gtt := 0
	for i := 0; i < n; i++ {
		resp := <-rchan
		if resp.err != nil {
			t.Error(resp.err)
			break
		}
		ctt += resp.connectTime
		gtt += resp.getSetTime
	}
	t.Logf("Avg connect time: %d Avg getset time: %d\n",
		int(ctt/n),
		int(gtt/n))
}

func testSetGetKey(t *testing.T, rchan chan resp, n int) {
	tu := time.Microsecond
	resp := resp{}
	st := time.Now()
	rw, err := connect()
	if err != nil {
		resp.err = err
		rchan <- resp
		return
	}
	resp.connectTime = int(time.Since(st) / tu)
	ns := str(n)

	st = time.Now()
	write(t, rw, "SET", "Lewis "+ns, "Hamilton"+ns)
	res := read(t, rw)
	fgst := int(time.Since(st) / tu)

	assert.Equal(t, "OK", res)

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	st = time.Now()
	write(t, rw, "GET", "Lewis "+ns)
	assert.Equal(t, "Hamilton"+ns, read(t, rw))
	sgst := int(time.Since(st) / tu)
	resp.getSetTime = int((fgst + sgst) / 2)

	rchan <- resp
}

func read(t *testing.T, r *bufio.ReadWriter) string {
	s, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("Exiting due to error", err)
		t.Fatalf("Read error %v", err)
	}

	if s[0] == '$' {
		s, err = r.ReadString('\n')
		if err != nil {
			fmt.Println("Exiting due to error", err)
			t.Fatalf("Read error %v", err)
		}
	} else {
		s = s[1:]
	}
	s = s[:len(s)-2]
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

func connect() (*bufio.ReadWriter, error) {
	conn, err := net.Dial("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Exiting due to error", err)
		return nil, err
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}
