package parser

import (
	"encoding/json"
)

type jsonParser struct {
	message *Message
}

func NewJsonParser() Parser {
	message := new(Message)
	return &jsonParser{
		message: message,
	}
}

func (p *jsonParser) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, p.message); err != nil {
		return err
	}
	return nil
}

func (p *jsonParser) Query() string {
	return p.message.Query
}
