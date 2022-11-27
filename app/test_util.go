package main

import (
	"redis-go/app/mocks"

	gomock "github.com/golang/mock/gomock"
)

func mockReadString(mr *mocks.MockStringReader, str string, err error) {
	mr.EXPECT().
		ReadString(gomock.Eq(byte('\n'))).
		Return(str, err)
}
