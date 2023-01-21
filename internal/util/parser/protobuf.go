package parser

import (
	"github.com/Team-OurPlayground/idl/goproto"
	"google.golang.org/protobuf/proto"
)

type protobufParser struct {
	protoData *goproto.Message
}

func NewProtobufParser() Parser {
	protoData := new(goproto.Message)
	return &protobufParser{protoData: protoData}
}

func (p *protobufParser) Unmarshal(data []byte) error {
	if err := proto.Unmarshal(data, p.protoData); err != nil {
		return err
	}
	return nil
}

func (p *protobufParser) Query() string {
	return p.protoData.Query
}
