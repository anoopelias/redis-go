package main

import (
	"fmt"
)

type Command interface {
	ReadParams(len int) error
	Execute() string
}

type CommandReader struct {
	respReader RespReader
}

func NewCommandReader(rr RespReader) CommandReader {
	return CommandReader{
		respReader: rr,
	}
}

func (cr *CommandReader) Read() (Command, error) {
	l, err := cr.respReader.ReadArrayLen()
	if err != nil {
		return nil, err
	}
	if l <= 0 {
		return nil, fmt.Errorf("array length should be at least 1")
	}

	cs, err := cr.respReader.ReadBulkString()
	if err != nil {
		return nil, err
	}

	var c Command
	switch cs {
	case "PING":
		c = NewPingCommand()
	case "ECHO":
		c = NewEchoCommand(cr.respReader)
	}

	err = c.ReadParams(l - 1)
	if err != nil {
		return nil, err
	}

	return c, nil
}

type PingCommand struct {
}

func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

func (c *PingCommand) ReadParams(len int) error {
	if len != 0 {
		return fmt.Errorf("invalid number of params")
	}
	return nil
}

func (c *PingCommand) Execute() string {
	return "+PONG"
}

type EchoCommand struct {
	reader RespReader
	str    string
}

func NewEchoCommand(rr RespReader) *EchoCommand {
	return &EchoCommand{
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

func (c *EchoCommand) Execute() string {
	return "+" + c.str
}
