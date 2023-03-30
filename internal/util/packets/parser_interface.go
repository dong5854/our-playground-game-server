package parser

import "github.com/Team-OurPlayground/idl/FBPacket"

type ChatParser interface {
	Parse(data []byte)
	Message() string
	SenderID() string
	ReceiverID() string
	Type() string
}

type Chat interface {
	Create(message string, senderID string, receiverID string, chatType FBPacket.ChatType) []byte
}

type Parser interface {
	Unmarshal(data []byte) error
	Function() string
	Data() string
	Dx() float32
	Dy() float32
}
