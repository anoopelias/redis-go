package redis_go

import (
	"redis-go/app/mocks"
	"strconv"

	gomock "github.com/golang/mock/gomock"
)

func mockReadString(mr *mocks.MockStringReader, str string, err error) {
	mr.EXPECT().
		ReadString(gomock.Eq(byte('\n'))).
		Return(str, err)
}

func mockReadCommand(mr *mocks.MockStringReader, err error,
	l int, str ...string) {
	mockReadString(mr, "*"+strconv.Itoa(l)+"\r\n", nil)
	for i, v := range str {
		mockReadString(mr, "$"+strconv.Itoa(len(v))+"\r\n", nil)
		if i == len(str)-1 && err != nil {
			mockReadString(mr, v+"\r\n", err)
		} else {
			mockReadString(mr, v+"\r\n", nil)
		}
	}
}
