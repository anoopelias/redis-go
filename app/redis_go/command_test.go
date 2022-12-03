package redis_go

import (
	"fmt"
	"redis-go/app/mocks"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCommandSetAndGetWithoutPx(t *testing.T) {
	data := make(map[string]string)
	sc := SetCommand{
		key:   "hello",
		value: "world",
		px:    -1,
	}
	sc.Execute(&data)
	gc := GetCommand{
		key: "hello",
	}
	assert.Equal(t, "+world", gc.Execute(&data))
	time.Sleep(600 * time.Millisecond)
	assert.Equal(t, "+world", gc.Execute(&data))
}

func TestCommandSetAndGetWithPx(t *testing.T) {
	data := make(map[string]string)
	sc := SetCommand{
		key:   "hello",
		value: "world",
		px:    500,
	}
	sc.Execute(&data)
	gc := GetCommand{
		key: "hello",
	}
	assert.Equal(t, "+world", gc.Execute(&data))
	time.Sleep(600 * time.Millisecond)
	assert.Equal(t, "$-1", gc.Execute(&data))
}

func TestCommandSetAndGet(t *testing.T) {
	data := make(map[string]string)
	sc := SetCommand{
		key:   "hello",
		value: "world",
	}
	sc.Execute(&data)
	gc := GetCommand{
		key: "hello",
	}
	assert.Equal(t, "+world", gc.Execute(&data))
}

func TestCommandGet(t *testing.T) {
	data := make(map[string]string)
	gc := GetCommand{}
	assert.Equal(t, gc.Execute(&data), "$-1")
}

func TestGetReadParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrr := mocks.NewMockRespReader(ctrl)

	ec := NewGetCommand(mrr)
	mrr.EXPECT().ReadBulkString().Return("mykey", nil)
	assert.Nil(t, ec.ReadParams(1))
	assert.Equal(t, ec.key, "mykey")
}

func TestGetReadParamsLenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrr := mocks.NewMockRespReader(ctrl)

	ec := NewGetCommand(mrr)
	assert.NotNil(t, ec.ReadParams(2))
}

func TestGetReadParamsKeyReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrr := mocks.NewMockRespReader(ctrl)

	ec := NewGetCommand(mrr)
	mrr.EXPECT().ReadBulkString().Return("", fmt.Errorf("read error"))
	assert.NotNil(t, ec.ReadParams(1))
}

func TestCommandSet(t *testing.T) {
	data := make(map[string]string)
	sc := SetCommand{}
	assert.Equal(t, sc.Execute(&data), "+OK")
}

func TestSetReadParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrr := mocks.NewMockRespReader(ctrl)

	ec := NewSetCommand(mrr)
	mrr.EXPECT().ReadBulkString().Return("mykey", nil)
	mrr.EXPECT().ReadBulkString().Return("myvalue", nil)
	assert.Nil(t, ec.ReadParams(2))
	assert.Equal(t, ec.key, "mykey")
	assert.Equal(t, ec.value, "myvalue")
}

func TestSetReadParamsLenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrr := mocks.NewMockRespReader(ctrl)

	sc := NewSetCommand(mrr)
	assert.NotNil(t, sc.ReadParams(0))
	mrr.EXPECT().ReadBulkString().Times(0)
}

func TestSetReadParamsKeyReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrr := mocks.NewMockRespReader(ctrl)

	ec := NewSetCommand(mrr)
	mrr.EXPECT().ReadBulkString().Return("Hello", fmt.Errorf("read error"))
	assert.NotNil(t, ec.ReadParams(2))
}

func TestSetReadParamsValueReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrr := mocks.NewMockRespReader(ctrl)

	ec := NewSetCommand(mrr)
	mrr.EXPECT().ReadBulkString().Return("Hello", nil)
	mrr.EXPECT().ReadBulkString().Return("World", fmt.Errorf("read error"))
	assert.NotNil(t, ec.ReadParams(2))
}

func TestCommandEchoExecute(t *testing.T) {
	data := make(map[string]string)
	ec := EchoCommand{str: "Hello World!"}
	assert.Equal(t, ec.Execute(&data), "+Hello World!")
}

func TestCommandPingExecute(t *testing.T) {
	data := make(map[string]string)
	pc := PingCommand{}
	assert.Equal(t, pc.Execute(&data), "+PONG")
}

func TestCommandReaderGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, nil, 2, "GET", "Hello")

	c, err := cr.Read()
	assert.Nil(t, err)

	ec := c.(*GetCommand)
	assert.NotNil(t, ec)
	assert.Equal(t, ec.key, "Hello")
}

func TestCommandReaderGetSmallCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, nil, 2, "get", "Hello")

	c, err := cr.Read()
	assert.Nil(t, err)

	ec := c.(*GetCommand)
	assert.NotNil(t, ec)
	assert.Equal(t, ec.key, "Hello")
}

func TestCommandReaderSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, nil, 3, "SET", "Hello", "World")

	c, err := cr.Read()
	assert.Nil(t, err)

	ec := c.(*SetCommand)
	assert.NotNil(t, ec)
	assert.Equal(t, ec.key, "Hello")
	assert.Equal(t, ec.value, "World")
	assert.Equal(t, ec.px, -1)
}

func TestCommandReaderSmallCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, nil, 3, "set", "Hello", "World")

	c, err := cr.Read()
	assert.Nil(t, err)

	ec := c.(*SetCommand)
	assert.NotNil(t, ec)
	assert.Equal(t, ec.key, "Hello")
	assert.Equal(t, ec.value, "World")
}

func TestCommandReaderEcho(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, nil, 2, "ECHO", "Hello World!")

	c, err := cr.Read()
	assert.Equal(t, err, nil)

	ec := c.(*EchoCommand)
	assert.NotNil(t, ec)
	assert.Equal(t, ec.str, "Hello World!")
}

func TestCommandReaderEchoLowerCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, nil, 2, "echo", "Hello World!")

	c, err := cr.Read()
	assert.Equal(t, err, nil)

	ec := c.(*EchoCommand)
	assert.NotNil(t, ec)
	assert.Equal(t, ec.str, "Hello World!")
}

func TestCommandReaderEchoArrayLenReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadString(mr, "*2\r\n", fmt.Errorf("read error"))

	_, err := cr.Read()
	assert.NotNil(t, err)
}

func TestCommandReaderEchoArrayLenZeroError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadString(mr, "*0\r\n", nil)

	_, err := cr.Read()
	assert.NotNil(t, err)
}

func TestCommandReaderEchoBulkStringReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, fmt.Errorf("read error"), 2, "ECHO")

	_, err := cr.Read()
	assert.NotNil(t, err)
}

func TestCommandReaderEchoReadParamsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, nil, 1, "ECHO")

	_, err := cr.Read()
	assert.NotNil(t, err)
}

func TestCommandReaderPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, nil, 1, "PING")

	c, err := cr.Read()
	assert.Equal(t, err, nil)

	pc := c.(*PingCommand)
	assert.NotNil(t, pc)
}

func TestCommandReaderPingLowerCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, nil, 1, "ping")

	c, err := cr.Read()
	assert.Equal(t, err, nil)

	pc := c.(*PingCommand)
	assert.NotNil(t, pc)
}

func TestCommandReaderSetWithPx(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, nil, 5, "set", "Lewis", "Hamilton", "PX", "100")

	c, err := cr.Read()
	assert.Equal(t, nil, err)

	ec := c.(*SetCommand)
	assert.NotNil(t, ec)
	assert.Equal(t, ec.key, "Lewis")
	assert.Equal(t, ec.value, "Hamilton")
	assert.Equal(t, ec.px, 100)
}

func TestCommandReaderSetWithPxSmallCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, nil, 5, "set", "Lewis", "Hamilton", "px", "100")

	c, err := cr.Read()
	assert.Equal(t, nil, err)

	ec := c.(*SetCommand)
	assert.NotNil(t, ec)
	assert.Equal(t, ec.key, "Lewis")
	assert.Equal(t, ec.value, "Hamilton")
	assert.Equal(t, ec.px, 100)
}

func TestCommandReaderSetWithUnknownParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, fmt.Errorf("read error"), 5,
		"set", "Lewis", "Hamilton", "gx")

	_, err := cr.Read()
	assert.NotNil(t, err)
}

func TestCommandReaderSetReadPxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, fmt.Errorf("read error"), 5,
		"set", "Lewis", "Hamilton", "px")

	_, err := cr.Read()
	assert.NotNil(t, err)
}

func TestCommandReaderSetPxReadTimeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, fmt.Errorf("read error"), 5,
		"set", "Lewis", "Hamilton", "px", "100")

	_, err := cr.Read()
	assert.NotNil(t, err)
}

func TestCommandReaderSetPxReadTimeStringError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadCommand(mr, nil, 5,
		"set", "Lewis", "Hamilton", "px", "abcd")

	_, err := cr.Read()
	assert.NotNil(t, err)
}
func TestCommandReaderUnknownCommandError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockStringReader(ctrl)
	rr := NewRespReader(mr)
	cr := NewCommandReader(rr)

	mockReadString(mr, "*1\r\n", nil)
	mockReadString(mr, "$4\r\n", nil)
	mockReadString(mr, "PICK\r\n", nil)

	_, err := cr.Read()
	assert.NotNil(t, err)
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
	mrr := mocks.NewMockRespReader(ctrl)

	ec := NewEchoCommand(mrr)
	mrr.EXPECT().ReadBulkString().Return("Hello World!", nil)
	assert.Nil(t, ec.ReadParams(1))
	assert.Equal(t, ec.str, "Hello World!")
}

func TestEchoReadParamsReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrr := mocks.NewMockRespReader(ctrl)

	ec := NewEchoCommand(mrr)
	mrr.EXPECT().ReadBulkString().Return("Hello", fmt.Errorf("read error"))
	assert.NotNil(t, ec.ReadParams(1))
}

func TestEchoReadParamsLenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mrr := mocks.NewMockRespReader(ctrl)

	pc := NewEchoCommand(mrr)
	assert.NotNil(t, pc.ReadParams(0))
	mrr.EXPECT().ReadBulkString().Times(0)
}
