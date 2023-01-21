package parser

import (
	idl "github.com/Team-OurPlayground/idl/proto"
	"google.golang.org/protobuf/proto"
)

type protobufParser struct {
	protoData *idl.SearchRequest
}

func NewProtobufParser() Parser {
	protoData := new(idl.SearchRequest)
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
