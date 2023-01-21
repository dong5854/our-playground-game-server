package parser

import (
	idl "github.com/Team-OurPlayground/idl/proto"
	"github.com/vmihailenco/msgpack"
)

type msgPackParser struct {
	message Message
}

func NewMsgPackParser() Parser {
	protoData := new(idl.SearchRequest)
	return &protobufParser{protoData: protoData}
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
