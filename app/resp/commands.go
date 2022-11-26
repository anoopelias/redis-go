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
	reader StringReader
}

func NewEchoCommand(sr StringReader) EchoCommand {
	return EchoCommand{
		reader: sr,
	}
}

func (c *EchoCommand) ReadParams(len int) error {
	if len != 1 {
		return fmt.Errorf("incorrect number of params")
	}

	return nil
}
