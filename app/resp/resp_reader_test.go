package resp

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReadBulkString(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("$12\r\n", nil)
	mr.mockReadString("Hello World!\r\n", nil)

	str, err := rr.ReadBulkString()
	assert.Nil(t, err)
	assert.Equal(t, str, "Hello World!")
}

func TestReadBulkStringLenReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("$12\r\n", fmt.Errorf("read error"))

	_, err := rr.ReadBulkString()
	assert.NotNil(t, err)
}

func TestReadBulkStringLenParseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("$abc\r\n", nil)

	_, err := rr.ReadBulkString()
	assert.NotNil(t, err)
}

func TestReadBulkStringLenTypeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("*12\r\n", nil)

	_, err := rr.ReadBulkString()
	assert.NotNil(t, err)
}

func TestReadBulkStringReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("$12\r\n", nil)
	mr.mockReadString("Hello World!\r\n", fmt.Errorf("read error"))

	_, err := rr.ReadBulkString()
	assert.NotNil(t, err)
}

func TestReadBulkStringLenMistmatchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("$13\r\n", nil)
	mr.mockReadString("Hello World!\r\n", nil)

	_, err := rr.ReadBulkString()
	assert.NotNil(t, err)
}

func TestReadLine(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("PING\r\n", nil)

	line, err := rr.ReadLine()

	assert.Equal(t, err, nil)
	assert.Equal(t, line, "PING")
}

func TestReadLineReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.mockReadString("PING\r\n", fmt.Errorf("read error"))

	_, err := rr.ReadLine()

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