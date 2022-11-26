package resp

import (
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
	mr := NewMockStringReader(ctrl)

	ec := NewEchoCommand(mr)
	// mr.mockReadString("$12\r\n", nil)
	// mr.mockReadString("Hello World!\r\n", nil)
	assert.Nil(t, ec.ReadParams(1))
}

func TestEchoReadParamsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mr := NewMockStringReader(ctrl)

	pc := NewEchoCommand(mr)
	assert.NotNil(t, pc.ReadParams(0))
	mr.EXPECT().ReadString(gomock.Any()).Times(0)
}
