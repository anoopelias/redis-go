package ev

import (
	"fmt"
	"redis-go/app/mocks"
	"syscall"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	var sa syscall.SockaddrInet4
	sa.Port = 6379
	sa.Addr = [4]byte{0, 0, 0, 0}

	sc.EXPECT().Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0).Return(254, nil)
	sc.EXPECT().Bind(254, &sa).Return(nil)
	sc.EXPECT().Listen(254, 50).Return(nil)
	sc.EXPECT().SetNonblock(254, true).Return(nil)

	sfd, err := el.create()
	assert.Nil(t, err)
	assert.True(t, sfd > 0)
}

func TestCreateError_Socket(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)

	sc.EXPECT().Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0).Return(254, fmt.Errorf("socket error"))

	_, err := el.create()
	assert.NotNil(t, err)
}

func TestCreateError_Bind(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	var sa syscall.SockaddrInet4
	sa.Port = 6379
	sa.Addr = [4]byte{0, 0, 0, 0}

	sc.EXPECT().Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0).Return(254, nil)
	sc.EXPECT().Bind(254, &sa).Return(fmt.Errorf("bind error"))

	_, err := el.create()
	assert.NotNil(t, err)
}

func TestCreateError_Listen(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	var sa syscall.SockaddrInet4
	sa.Port = 6379
	sa.Addr = [4]byte{0, 0, 0, 0}

	sc.EXPECT().Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0).Return(254, nil)
	sc.EXPECT().Bind(254, &sa).Return(nil)
	sc.EXPECT().Listen(254, 50).Return(fmt.Errorf("listen error"))

	_, err := el.create()
	assert.NotNil(t, err)
}

func TestCreateError_SetNonblock(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	var sa syscall.SockaddrInet4
	sa.Port = 6379
	sa.Addr = [4]byte{0, 0, 0, 0}

	sc.EXPECT().Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0).Return(254, nil)
	sc.EXPECT().Bind(254, &sa).Return(nil)
	sc.EXPECT().Listen(254, 50).Return(nil)
	sc.EXPECT().SetNonblock(254, true).Return(fmt.Errorf("non-block error"))

	_, err := el.create()
	assert.NotNil(t, err)
}

func TestExecute(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)
	sce := mocks.NewMockSysCallError(ctrl)

	el := NewSocketEventLoop(sc)
	el.cfds = []int{495}
	sce.EXPECT().Timeout().Return(true).Times(2)
	sc.EXPECT().Accept(245).Return(-1, nil, sce)
	sc.EXPECT().Read(495, gomock.Any()).Return(-1, sce)

	err := el.execute(245)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(el.cfds))
	assert.Equal(t, 495, el.cfds[0])
}

func TestExecuteDisconnected(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)
	sce := mocks.NewMockSysCallError(ctrl)

	el := NewSocketEventLoop(sc)
	el.cfds = []int{495}
	sce.EXPECT().Timeout().Return(true)
	sc.EXPECT().Accept(245).Return(-1, nil, sce)
	sc.EXPECT().Read(495, gomock.Any()).Return(0, nil)
	sc.EXPECT().Close(495).Return(nil)

	err := el.execute(245)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(el.cfds))
}

func TestExecuteDisconnected_CloseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)
	sce := mocks.NewMockSysCallError(ctrl)

	el := NewSocketEventLoop(sc)
	el.cfds = []int{495}
	sce.EXPECT().Timeout().Return(true)
	sc.EXPECT().Accept(245).Return(-1, nil, sce)
	sc.EXPECT().Read(495, gomock.Any()).Return(0, nil)
	sc.EXPECT().Close(495).Return(fmt.Errorf("close error"))

	err := el.execute(245)
	assert.NotNil(t, err)
}

func TestExecuteExistingCfd(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	sc.EXPECT().Accept(245).Return(455, nil, nil)
	sc.EXPECT().SetNonblock(455, true).Return(nil)
	sc.EXPECT().Read(455, gomock.Any()).Return(3, nil)
	el.handler = func(sr StringReader) string {
		return "Ok"
	}
	sc.EXPECT().Write(455, gomock.Any()).Return(2, nil)

	err := el.execute(245)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(el.cfds))
	assert.Equal(t, 455, el.cfds[0])
}

func TestExecuteError_AcceptEagain(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)
	scerr := mocks.NewMockSysCallError(ctrl)

	el := NewSocketEventLoop(sc)
	sc.EXPECT().Accept(245).Return(-1, nil, scerr)
	scerr.EXPECT().Timeout().Return(true)

	err := el.execute(245)
	assert.Nil(t, err)
}

func TestExecuteError_Accept(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)
	scerr := mocks.NewMockSysCallError(ctrl)

	el := NewSocketEventLoop(sc)
	sc.EXPECT().Accept(245).Return(455, nil, scerr)
	scerr.EXPECT().Timeout().Return(false)

	err := el.execute(245)
	assert.NotNil(t, err)
}

func TestExecuteError_AcceptNotTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	sc.EXPECT().Accept(245).Return(455, nil, fmt.Errorf("accept error"))

	err := el.execute(245)
	assert.NotNil(t, err)
}

func TestExecuteError_SetNonblock(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	sc.EXPECT().Accept(245).Return(455, nil, nil)
	sc.EXPECT().SetNonblock(455, true).Return(fmt.Errorf("nonblock error"))

	err := el.execute(245)
	assert.NotNil(t, err)
}

func TestRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	sc.EXPECT().Read(455, gomock.Any()).DoAndReturn(func(_ int, data []byte) (int, error) {
		data[0] = 67
		data[1] = 104
		data[2] = 10
		return 3, nil
	})

	data, err := el.read(455)
	assert.Nil(t, err)
	assert.Equal(t, []byte{67, 104, 10}, data)
}

func TestReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	sc.EXPECT().Read(455, gomock.Any()).Return(0, fmt.Errorf("read error"))

	_, err := el.read(455)
	assert.NotNil(t, err)
}

func TestProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	sc.EXPECT().Read(455, gomock.Any()).DoAndReturn(func(_ int, data []byte) (int, error) {
		data[0] = 67
		data[1] = 104
		data[2] = 10
		return 3, nil
	})

	el.handler = func(sr StringReader) string {
		str, err := sr.ReadString('\n')
		assert.Nil(t, err)
		assert.Equal(t, "Ch\n", str)

		_, err = sr.ReadString('\n')
		assert.NotNil(t, err)
		return "+OK"
	}
	res := []byte("+OK")
	sc.EXPECT().Write(455, res).Return(len(res), nil)

	ctd, err := el.process(455)
	assert.Nil(t, err)
	assert.Equal(t, true, ctd)
}

func TestProcessError_Read(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	sc.EXPECT().Read(455, gomock.Any()).Return(-1, fmt.Errorf("read error"))

	_, err := el.process(455)
	assert.NotNil(t, err)
}

func TestProcessError_ReadRetry(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)
	sce := mocks.NewMockSysCallError(ctrl)

	el := NewSocketEventLoop(sc)
	sce.EXPECT().Timeout().Return(true)
	sc.EXPECT().Read(455, gomock.Any()).Return(-1, sce)

	ctd, err := el.process(455)
	assert.Nil(t, err)
	assert.True(t, ctd)
}

func TestProcessError_Disconnected(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	sc.EXPECT().Read(455, gomock.Any()).Return(0, nil)

	ctd, err := el.process(455)
	assert.Nil(t, err)
	assert.False(t, ctd)
}

func TestProcessError_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	sc.EXPECT().Read(455, gomock.Any()).Return(3, nil)

	el.handler = func(sr StringReader) string {
		return "Ok"
	}
	sc.EXPECT().Write(455, gomock.Any()).Return(3, fmt.Errorf("write error"))

	_, err := el.process(455)
	assert.NotNil(t, err)
}

func TestProcessError_WriteIncomplete(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	sc.EXPECT().Read(455, gomock.Any()).Return(3, nil)

	el.handler = func(sr StringReader) string {
		return "Ok"
	}
	sc.EXPECT().Write(455, gomock.Any()).Return(1, nil)

	_, err := el.process(455)
	assert.NotNil(t, err)
}
