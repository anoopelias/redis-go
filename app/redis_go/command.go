package redis_go

import (
	"fmt"
)

type Command interface {
	ReadParams(len int) error
	Execute(data *map[string]string) string
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
	case "PING", "ping":
		c = NewPingCommand()
	case "ECHO", "echo":
		c = NewEchoCommand(cr.respReader)
	default:
		return nil, fmt.Errorf("unknown command %s", cs)
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

func (c *PingCommand) Execute(data *map[string]string) string {
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

func (c *EchoCommand) Execute(data *map[string]string) string {
	return "+" + c.str
}

type SetCommand struct {
	reader RespReader
	key    string
	value  string
}

func NewSetCommand(rr RespReader) *SetCommand {
	return &SetCommand{
		reader: rr,
	}
}

func (s *SetCommand) ReadParams(len int) (err error) {
	if len != 2 {
		return fmt.Errorf("incorrect number of params")
	}
	key, err := s.reader.ReadBulkString()
	if err != nil {
		return
	}

	value, err := s.reader.ReadBulkString()
	if err != nil {
		return
	}

	s.key = key
	s.value = value
	return nil
}

func (s *SetCommand) Execute(data *map[string]string) string {
	(*data)[s.key] = s.value
	return "+OK"
}

type GetCommand struct {
	reader RespReader
	key    string
}

func (g *GetCommand) ReadParams(len int) (err error) {
	if len != 1 {
		return fmt.Errorf("incorrect number of params")
	}
	key, err := g.reader.ReadBulkString()
	if err != nil {
		return
	}

	g.key = key
	return nil
}

func NewGetCommand(rr RespReader) *GetCommand {
	return &GetCommand{
		reader: rr,
	}
}

func (g *GetCommand) Execute(data *map[string]string) string {
	if val, ok := (*data)[g.key]; ok {
		return "+" + val
	}
	return "$-1"
}
