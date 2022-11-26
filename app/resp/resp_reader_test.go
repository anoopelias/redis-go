package resp

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReadCommandEcho(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("*2\r\n", nil)
	mr.mockReadString("$4\r\n", nil)
	mr.mockReadString("ECHO\r\n", nil)
	// mr.mockReadString("$12\r\n", nil)
	// mr.mockReadString("Hello World!\r\n", nil)

	_, err := rr.ReadCommand()
	assert.Equal(t, err, nil)

	// pc := c.(*EchoCommand)
	// assert.NotNil(t, pc)
}

func TestReadCommandPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("*1\r\n", nil)
	mr.mockReadString("$4\r\n", nil)
	mr.mockReadString("PING\r\n", nil)

	c, err := rr.ReadCommand()
	assert.Equal(t, err, nil)

	pc := c.(*PingCommand)
	assert.NotNil(t, pc)
}

func TestReadCommandsArrayErrorFirstTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("*1\r\n", fmt.Errorf(""))

	_, err := rr.ReadCommand()
	assert.NotEqual(t, err, nil)
}

func TestReadCommandsArrayErrorSecondTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("*1\r\n", nil)
	mr.mockReadString("*1\r\n", fmt.Errorf(""))

	_, err := rr.ReadCommand()
	assert.NotEqual(t, err, nil)
}

func TestReadLine(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("PING\r\n", nil)

	line, err := rr.readLine()

	assert.Equal(t, err, nil)
	assert.Equal(t, line, "PING")
}

func TestReadLineReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("PING\r\n", fmt.Errorf("read error"))

	_, err := rr.readLine()

	assert.NotEqual(t, err, nil)
}

func TestParseInt(t *testing.T) {
	ty, d, err := parse(":23")

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if ty != respInt {
		t.Errorf("Incorrect type")
	}

	if d != 23 {
		t.Errorf("Incorrect value %v", d)
	}
}

func TestParseIntError(t *testing.T) {
	_, _, err := parse(":1a")

	if err == nil {
		t.Errorf("Unexpected error")
	}

}

func TestParseIntTrim(t *testing.T) {
	ty, d, err := parse(" :23 ")

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if ty != respInt {
		t.Errorf("Incorrect type")
	}

	if d != 23 {
		t.Errorf("Incorrect value %v", d)
	}
}

func TestParseString(t *testing.T) {
	ty, d, err := parse("+Hotel")

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if ty != respString {
		t.Errorf("Incorrect type")
	}

	if d != "Hotel" {
		t.Errorf("Incorrect value %v", d)
	}
}

func TestParseArray(t *testing.T) {
	ty, d, err := parse("*5")

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if ty != respArrayLen {
		t.Errorf("Incorrect type")
	}

	if d != 5 {
		t.Errorf("Incorrect value %v", d)
	}
}

func TestParseArrayError(t *testing.T) {
	_, _, err := parse("*1a")

	if err == nil {
		t.Errorf("Unexpected error")
	}

}

func TestParseBulkStringLen(t *testing.T) {
	ty, d, err := parse("$23")

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if ty != respBulkStringLen {
		t.Errorf("Incorrect type")
	}

	if d != 23 {
		t.Errorf("Incorrect value %v", d)
	}
}

func TestParseBulkStringLenError(t *testing.T) {
	_, _, err := parse("$1a")

	if err == nil {
		t.Errorf("Unexpected error")
	}

}
