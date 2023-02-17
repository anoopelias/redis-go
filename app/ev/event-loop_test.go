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
	sc.EXPECT().Kqueue().Return(375, nil)
	sc.EXPECT().Kevent(375, eventsFor(254), nil, nil).Return(0, nil)

	err := el.create()
	assert.Nil(t, err)
	assert.True(t, el.sfd > 0)
	assert.True(t, el.kq > 0)
}

func TestCreateError_Socket(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)

	sc.EXPECT().Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0).Return(254, fmt.Errorf("socket error"))

	err := el.create()
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

	err := el.create()
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

	err := el.create()
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

	err := el.create()
	assert.NotNil(t, err)
}

func TestCreateError_Kqueue(t *testing.T) {
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
	sc.EXPECT().Kqueue().Return(375, fmt.Errorf("kqueue creation failed"))

	err := el.create()
	assert.NotNil(t, err)
}

func TestCreateError_Kevent(t *testing.T) {
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
	sc.EXPECT().Kqueue().Return(375, nil)
	sc.EXPECT().Kevent(375, eventsFor(254), nil, nil).Return(0, fmt.Errorf("kevent error"))

	err := el.create()
	assert.NotNil(t, err)
}

func TestExecuteNoEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	el.sfd = 245
	el.kq = 375

	events := make([]syscall.Kevent_t, 10)
	sc.EXPECT().Kevent(375, nil, events, nil).Return(0, nil)

	err := el.execute()
	assert.Nil(t, err)
}

func TestExecuteError_Temporary(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)
	sce := mocks.NewMockSysCallError(ctrl)

	el := NewSocketEventLoop(sc)
	el.sfd = 245
	el.kq = 375

	events := make([]syscall.Kevent_t, 10)
	sce.EXPECT().Temporary().Return(true)
	sc.EXPECT().Kevent(375, nil, events, nil).Return(0, sce)

	err := el.execute()
	assert.Nil(t, err)
}

func TestExecuteError_FetchEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	el.sfd = 245
	el.kq = 375

	events := make([]syscall.Kevent_t, 10)
	sc.EXPECT().Kevent(375, nil, events, nil).Return(0, fmt.Errorf("fetch events"))

	err := el.execute()
	assert.NotNil(t, err)
}

func TestExecuteSfdEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	el.sfd = 245
	el.kq = 375

	events := make([]syscall.Kevent_t, 10)
	sc.EXPECT().Kevent(375, nil, events, nil).DoAndReturn(funcKevent(245))

	sc.EXPECT().Accept(245).Return(455, nil, nil)
	sc.EXPECT().SetNonblock(455, true).Return(nil)
	sc.EXPECT().Kevent(375, eventsFor(455), nil, nil).Return(0, nil)

	err := el.execute()
	assert.Nil(t, err)
}

func TestExecuteError_Kevent(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	el.sfd = 245
	el.kq = 375

	events := make([]syscall.Kevent_t, 10)
	sc.EXPECT().Kevent(375, nil, events, nil).DoAndReturn(funcKevent(245))

	sc.EXPECT().Accept(245).Return(455, nil, nil)
	sc.EXPECT().SetNonblock(455, true).Return(nil)
	sc.EXPECT().Kevent(375, eventsFor(455), nil, nil).Return(0, fmt.Errorf("kev error"))

	err := el.execute()
	assert.NotNil(t, err)
}

func TestExecuteError_AcceptEagain(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)
	scerr := mocks.NewMockSysCallError(ctrl)

	el := NewSocketEventLoop(sc)
	events := make([]syscall.Kevent_t, 10)
	el.sfd = 245
	el.kq = 375

	sc.EXPECT().Kevent(375, nil, events, nil).DoAndReturn(funcKevent(245))
	sc.EXPECT().Accept(245).Return(-1, nil, scerr)
	scerr.EXPECT().Temporary().Return(true)

	err := el.execute()
	assert.Nil(t, err)
}

func TestExecuteError_Accept(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)
	scerr := mocks.NewMockSysCallError(ctrl)

	el := NewSocketEventLoop(sc)
	el.sfd = 245
	el.kq = 375

	events := make([]syscall.Kevent_t, 10)
	sc.EXPECT().Kevent(375, nil, events, nil).DoAndReturn(funcKevent(245))
	sc.EXPECT().Accept(245).Return(455, nil, scerr)
	scerr.EXPECT().Temporary().Return(false)

	err := el.execute()
	assert.NotNil(t, err)
}

func TestExecuteError_AcceptNotTemporary(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	el.sfd = 245
	el.kq = 375

	events := make([]syscall.Kevent_t, 10)
	sc.EXPECT().Kevent(375, nil, events, nil).DoAndReturn(funcKevent(245))
	sc.EXPECT().Accept(245).Return(455, nil, fmt.Errorf("accept error"))

	err := el.execute()
	assert.NotNil(t, err)
}

func TestExecuteError_NonBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	el.sfd = 245
	el.kq = 375

	events := make([]syscall.Kevent_t, 10)
	sc.EXPECT().Kevent(375, nil, events, nil).DoAndReturn(funcKevent(245))
	sc.EXPECT().Accept(245).Return(455, nil, nil)
	sc.EXPECT().SetNonblock(455, true).Return(fmt.Errorf("non-block error"))

	err := el.execute()
	assert.NotNil(t, err)
}

func TestExecuteCfd(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)
	sce := mocks.NewMockSysCallError(ctrl)

	el := NewSocketEventLoop(sc)
	el.sfd = 245
	el.kq = 375

	events := make([]syscall.Kevent_t, 10)
	sc.EXPECT().Kevent(375, nil, events, nil).DoAndReturn(funcKevent(495))
	sce.EXPECT().Temporary().Return(true)
	sc.EXPECT().Read(495, gomock.Any()).Return(-1, sce)

	err := el.execute()
	assert.Nil(t, err)
}

func TestExecuteCfdDisconnected(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	el.sfd = 245
	el.kq = 375

	events := make([]syscall.Kevent_t, 10)
	sc.EXPECT().Kevent(375, nil, events, nil).DoAndReturn(funcKevent(495))
	sc.EXPECT().Read(495, gomock.Any()).Return(0, nil)
	sc.EXPECT().Close(495).Return(nil)

	err := el.execute()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(el.cfds))
}

func TestExecuteDisconnected_CloseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	el.sfd = 245
	el.kq = 375

	events := make([]syscall.Kevent_t, 10)
	sc.EXPECT().Kevent(375, nil, events, nil).DoAndReturn(funcKevent(495))
	sc.EXPECT().Read(495, gomock.Any()).Return(0, nil)
	sc.EXPECT().Close(495).Return(fmt.Errorf("close error"))

	err := el.execute()
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
	sce.EXPECT().Temporary().Return(true)
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

func TestAddKqEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	el.kq = 300
	sc.EXPECT().Kevent(300, eventsFor(375), nil, nil).Return(0, nil)

	err := el.addKqEvent(375)
	assert.Nil(t, err)
}

func TestAddKqEventError(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := mocks.NewMockSysCall(ctrl)

	el := NewSocketEventLoop(sc)
	el.kq = 300
	sc.EXPECT().Kevent(300, gomock.Any(), nil, nil).Return(0, fmt.Errorf("kevent error"))

	err := el.addKqEvent(375)
	assert.NotNil(t, err)
}

func eventsFor(fd int) []syscall.Kevent_t {
	return []syscall.Kevent_t{{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD,
	}}
}

func funcKevent(ident int) func(int, []syscall.Kevent_t, []syscall.Kevent_t, *syscall.Timespec) (n int, err error) {
	return func(_ int, _ []syscall.Kevent_t, events []syscall.Kevent_t, _ *syscall.Timespec) (int, error) {
		events[0].Ident = uint64(ident)
		return 1, nil
	}
}
