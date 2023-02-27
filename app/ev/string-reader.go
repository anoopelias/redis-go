package ev

import "fmt"

type StringReader interface {
	ReadString(delim byte) (string, error)
}

type ArrayStringReader struct {
	pos int
	arr []byte
}

func NewArrayStringReader(arr []byte) StringReader {
	return &ArrayStringReader{
		arr: arr,
	}
}

func (a *ArrayStringReader) ReadString(delim byte) (string, error) {
	for i := a.pos; i < len(a.arr); i++ {
		if a.arr[i] == delim {
			str := string(a.arr[a.pos:(i + 1)])
			a.pos = i + 1
			return str, nil
		}
	}
	return "", fmt.Errorf("end of data")
}
