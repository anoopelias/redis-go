package main

import gomock "github.com/golang/mock/gomock"

func (mr *MockStringReader) mockReadString(str string, err error) {
	mr.EXPECT().
		ReadString(gomock.Eq(byte('\n'))).
		Return(str, err)
}
