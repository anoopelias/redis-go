package resp

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCommandReaderEcho(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mr.mockReadString("*2\r\n", nil)
	mr.mockReadString("$4\r\n", nil)
	mr.mockReadString("ECHO\r\n", nil)
	mr.mockReadString("$12\r\n", nil)
	mr.mockReadString("Hello World!\r\n", nil)

	c, err := cr.Read()
	assert.Equal(t, err, nil)

	ec := c.(*EchoCommand)
	assert.NotNil(t, ec)
	assert.Equal(t, ec.str, "Hello World!")
}

func TestCommandReaderEchoArrayLenReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mr.mockReadString("*2\r\n", fmt.Errorf("read error"))

	_, err := cr.Read()
	assert.NotNil(t, err)
}

func TestCommandReaderEchoArrayLenZeroError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mr.mockReadString("*0\r\n", nil)

	_, err := cr.Read()
	assert.NotNil(t, err)
}

func TestCommandReaderEchoBulkStringReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mr.mockReadString("*2\r\n", nil)
	mr.mockReadString("$4\r\n", nil)
	mr.mockReadString("ECHO\r\n", fmt.Errorf("read error"))

	_, err := cr.Read()
	assert.NotNil(t, err)
}

func TestCommandReaderEchoReadParamsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mr.mockReadString("*1\r\n", nil)
	mr.mockReadString("$4\r\n", nil)
	mr.mockReadString("ECHO\r\n", nil)

	_, err := cr.Read()
	assert.NotNil(t, err)
}

func TestCommandReaderPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mr.mockReadString("*1\r\n", nil)
	mr.mockReadString("$4\r\n", nil)
	mr.mockReadString("PING\r\n", nil)

	c, err := cr.Read()
	assert.Equal(t, err, nil)

	pc := c.(*PingCommand)
	assert.NotNil(t, pc)
}

func TestPingReadParams(t *testing.T) {
	pc := NewPingCommand()
	assert.Nil(t, pc.ReadParams(0))
}

func TestPingReadParamsError(t *testing.T) {
	pc := NewPingCommand()
	assert.NotNil(t, pc.ReadParams(1))
}

func TestEchoReadParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrr := NewMockRespReader(ctrl)

	ec := NewEchoCommand(mrr)
	mrr.EXPECT().ReadBulkString().Return("Hello World!", nil)
	assert.Nil(t, ec.ReadParams(1))
	assert.Equal(t, ec.str, "Hello World!")
}

func TestEchoReadParamsReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrr := NewMockRespReader(ctrl)

	ec := NewEchoCommand(mrr)
	mrr.EXPECT().ReadBulkString().Return("Hello", fmt.Errorf("read error"))
	assert.NotNil(t, ec.ReadParams(1))
}

func TestEchoReadParamsLenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrr := NewMockRespReader(ctrl)

	pc := NewEchoCommand(mrr)
	assert.NotNil(t, pc.ReadParams(0))
	mrr.EXPECT().ReadBulkString().Times(0)
}
