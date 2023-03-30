package packets

import (
	"github.com/Team-OurPlayground/idl/goproto"
	"google.golang.org/protobuf/proto"
)

type protobufParser struct {
	protoData *goproto.Data
}

func NewProtobufParser() Parser {
	protoData := new(goproto.Data)
	return &protobufParser{protoData: protoData}
}

func (p *protobufParser) Unmarshal(data []byte) error {
	if err := proto.Unmarshal(data, p.protoData); err != nil {
		return err
	}
	return nil
}

func (p *protobufParser) Function() string {
	return p.protoData.Function
}

func (p *protobufParser) Data() string {
	return p.protoData.Data
}

func (p *protobufParser) Dx() float32 {
	return p.protoData.Dx
}

func (p *protobufParser) Dy() float32 {
	return p.protoData.Dy
}
