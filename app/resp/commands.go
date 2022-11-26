package resp

import (
	"fmt"
)

type Command interface {
	ReadParams(len int) error
}

type PingCommand struct {
}

func NewPingCommand() PingCommand {
	return PingCommand{}
}

func (c *PingCommand) ReadParams(len int) error {
	if len != 0 {
		return fmt.Errorf("invalid number of params")
	}
	return nil
}

type EchoCommand struct {
	reader RespReader
	str    string
}

func NewEchoCommand(rr RespReader) EchoCommand {
	return EchoCommand{
		reader: rr,
	}
}

func (c *EchoCommand) ReadParams(len int) (err error) {
	if len != 1 {
		return fmt.Errorf("incorrect number of params")
	}
	str, err := c.reader.ReadBulkString()
	if err != nil {
		return
	}

	c.str = str
	return nil
}
