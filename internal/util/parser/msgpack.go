package parser

import (
	"github.com/vmihailenco/msgpack"
)

type msgPackParser struct {
	message *Message
}

func NewMsgPackParser() Parser {
	message := new(Message)
	return &msgPackParser{
		message: message,
	}
}

func (p *msgPackParser) Unmarshal(data []byte) error {
	if err := msgpack.Unmarshal(data, p.message); err != nil {
		return err
	}
	return nil
}

func (p *msgPackParser) Query() string {
	return p.message.Query
}
