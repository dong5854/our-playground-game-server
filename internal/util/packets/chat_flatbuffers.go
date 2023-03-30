package packets

import (
	"github.com/Team-OurPlayground/idl/FBPacket"
	flatbuffers "github.com/google/flatbuffers/go"
)

type chatPacket struct {
	chatData *FBPacket.ChatMessage
}

func NewChatParser() ChatParser {
	chatData := new(FBPacket.ChatMessage)
	return &chatPacket{chatData: chatData}
}

func NewChatCreator() ChatCreator {
	chatData := new(FBPacket.ChatMessage)
	return &chatPacket{chatData: chatData}
}

func (c *chatPacket) Create(message string, senderID string, receiverID string, chatType FBPacket.ChatType) []byte {
	builder := flatbuffers.NewBuilder(1024)

	msg := builder.CreateString(message)
	sender := builder.CreateString(senderID)
	receiver := builder.CreateString(receiverID)
	FBPacket.ChatMessageStart(builder)
	FBPacket.ChatMessageAddMessage(builder, msg)
	FBPacket.ChatMessageAddSender(builder, sender)
	FBPacket.ChatMessageAddReceiver(builder, receiver)
	FBPacket.ChatMessageAddType(builder, chatType)

	finalMsg := FBPacket.ChatMessageEnd(builder)
	builder.Finish(finalMsg)

	return builder.FinishedBytes()
}

func (c *chatPacket) Parse(data []byte) {
	c.chatData = FBPacket.GetRootAsChatMessage(data, 0)
}

func (c *chatPacket) Message() string {
	return string(c.chatData.Message())
}

func (c *chatPacket) SenderID() string {
	return string(c.chatData.Sender())
}

func (c *chatPacket) ReceiverID() string {
	return string(c.chatData.Receiver())
}

func (c *chatPacket) Type() string {
	return FBPacket.EnumNamesChatType[c.chatData.Type()]
}
