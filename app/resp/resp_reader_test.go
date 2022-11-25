package resp

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReadCommandsArraySuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.EXPECT().
		ReadString(gomock.Eq(byte('\n'))).
		Return("$4", nil).
		Times(3)

	err := rr.ReadCommand()
	assert.Equal(t, err, nil)
}

func TestReadCommandsArrayErrorFirstTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.EXPECT().
		ReadString(gomock.Eq(byte('\n'))).
		Return("$4", fmt.Errorf(""))

	err := rr.ReadCommand()
	assert.NotEqual(t, err, nil)
}

func TestReadCommandsArrayErrorSecondTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.EXPECT().
		ReadString(gomock.Eq(byte('\n'))).
		Return("$4", nil).
		Return("$4", fmt.Errorf(""))

	err := rr.ReadCommand()
	assert.NotEqual(t, err, nil)
}

func TestReadLineSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.EXPECT().
		ReadString(gomock.Eq(byte('\n'))).
		DoAndReturn(func(_ byte) (string, error) {
			return "$4", nil
		})

	ty, err := rr.readLine()

	assert.Equal(t, err, nil)
	assert.Equal(t, ty, respBulkStringLen)
}

func TestReadLineReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.EXPECT().
		ReadString(gomock.Eq(byte('\n'))).
		DoAndReturn(func(_ byte) (string, error) {
			return "$4", fmt.Errorf("read error")
		})

	_, err := rr.readLine()

	assert.NotEqual(t, err, nil)

}

func TestRespReaderReadLine_InvalidStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)

	mr.EXPECT().
		ReadString(gomock.Eq(byte('\n'))).
		DoAndReturn(func(_ byte) (string, error) {
			return "$4", nil
		})

	ty, err := rr.readLine()

	if err != nil {
		t.Errorf("No error for invalid array length, %v", err)
	}

	if ty != respBulkStringLen {
		t.Errorf("Incorrect type, %v", ty)
	}
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
