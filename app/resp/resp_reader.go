package resp

import (
	"fmt"
	"strconv"
	"strings"
)

type StringReader interface {
	ReadString(delim byte) (string, error)
}

type RespReader struct {
	reader StringReader
}

func NewRespReader(reader StringReader) RespReader {
	return RespReader{
		reader: reader,
	}
}

func (r *RespReader) ReadCommand() (err error) {

	for i := 0; i < 3; i++ {
		_, err = r.readLine()
		if err != nil {
			return
		}
	}

	return nil
}

func (r *RespReader) readLine() (t respType, err error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return
	}

	t, _, err = parse(line)
	return
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
