package resp

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
