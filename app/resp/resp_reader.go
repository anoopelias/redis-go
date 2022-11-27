package resp

import (
	"fmt"
	"strconv"
	"strings"
)

type StringReader interface {
	ReadString(delim byte) (string, error)
}

type RespReader interface {
	ReadLine() (string, error)
	ReadBulkString() (string, error)
	ReadArrayLen() (int, error)
}

type RespReaderImpl struct {
	reader StringReader
}

func NewRespReader(reader StringReader) RespReader {
	return &RespReaderImpl{
		reader: reader,
	}
}

func (r *RespReaderImpl) ReadArrayLen() (int, error) {
	line, err := r.ReadLine()
	if err != nil {
		return -1, err
	}

	t, v, err := parse(line)
	if err != nil {
		return -1, err
	}

	if t != respArrayLen {
		return -1, fmt.Errorf("expected array len %d", t)
	}

	return v.(int), nil
}

func (r *RespReaderImpl) ReadBulkString() (string, error) {
	line, err := r.ReadLine()
	if err != nil {
		return "", err
	}

	t, v, err := parse(line)
	if err != nil {
		return "", err
	}

	if t != respBulkStringLen {
		return "", fmt.Errorf("expected bulk string size %d", t)
	}

	line, err = r.ReadLine()
	if err != nil {
		return "", err
	}

	if len(line) != v.(int) {
		return "", fmt.Errorf("mismatched line length %d, %d", len(line), v.(int))
	}

	return line, nil
}

func (r *RespReaderImpl) ReadLine() (line string, err error) {
	line, err = r.reader.ReadString('\n')
	if err != nil {
		return
	}

	// remove \r\n from the end
	return line[:len(line)-2], nil
}

type respType int

const (
	invalid respType = iota
	respString
	respError
	respInt
	respArrayLen
	respBulkStringLen
)

func parse(s string) (respType, interface{}, error) {
	s = strings.Trim(s, " ")

	if s[0] == ':' || s[0] == '*' || s[0] == '$' {
		n, err := strconv.Atoi(s[1:])
		if err != nil {
			return invalid, nil, fmt.Errorf("cannot parse string to int")
		}

		if s[0] == ':' {
			return respInt, n, nil
		} else if s[0] == '*' {
			return respArrayLen, n, nil
		} else {
			return respBulkStringLen, n, nil
		}
	}

	if s[0] == '+' {
		return respString, s[1:], nil
	}

	return invalid, nil, fmt.Errorf("unknown type")
}
