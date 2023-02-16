package redis_go

import (
	"fmt"
	"strconv"
	"time"
)

type Command interface {
	ReadParams(len int) error
	Execute(data *map[string]string) string
	Response() chan string
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
	case "SET", "set":
		c = NewSetCommand(cr.respReader)
	case "GET", "get":
		c = NewGetCommand(cr.respReader)
	default:
		return nil, fmt.Errorf("unknown command %s", cs)
	}

	err = c.ReadParams(l - 1)
	if err != nil {
		return nil, err
	}

	return c, nil
}

type BaseCommand struct {
	resp chan string
}

func (c *BaseCommand) Response() chan string {
	return c.resp
}

func NewBaseCommand() BaseCommand {
	return BaseCommand{
		resp: make(chan string),
	}
}

type PingCommand struct {
	BaseCommand
}

func NewPingCommand() *PingCommand {
	return &PingCommand{
		BaseCommand: NewBaseCommand(),
	}
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
	BaseCommand
	reader RespReader
	str    string
}

func NewEchoCommand(rr RespReader) *EchoCommand {
	return &EchoCommand{
		BaseCommand: NewBaseCommand(),
		reader:      rr,
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
	BaseCommand
	reader RespReader
	key    string
	value  string
	px     int
}

func NewSetCommand(rr RespReader) *SetCommand {
	return &SetCommand{
		BaseCommand: NewBaseCommand(),
		reader:      rr,
		px:          -1,
	}
}

func (s *SetCommand) ReadParams(len int) (err error) {
	if len != 2 && len != 4 {
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

	if len == 2 {
		return
	}

	px, err := s.reader.ReadBulkString()
	if err != nil {
		return
	}

	if px != "PX" && px != "px" {
		return fmt.Errorf("expected px as the third arg")
	}

	ts, err := s.reader.ReadBulkString()
	if err != nil {
		return
	}

	t, err := strconv.Atoi(ts)
	if err != nil {
		return
	}

	s.px = t

	return nil
}

func (s *SetCommand) Execute(data *map[string]string) string {
	(*data)[s.key] = s.value
	if s.px != -1 {
		pxTimer := time.NewTimer(time.Duration(s.px) * time.Millisecond)
		go func() {
			<-pxTimer.C
			delete(*data, s.key)
		}()
	}
	return "+OK"
}

type GetCommand struct {
	BaseCommand
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
		BaseCommand: NewBaseCommand(),
		reader:      rr,
	}
}

func (g *GetCommand) Execute(data *map[string]string) string {
	if val, ok := (*data)[g.key]; ok {
		return "+" + val
	}
	return "$-1"
}
