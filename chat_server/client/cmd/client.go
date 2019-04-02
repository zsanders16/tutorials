package cmd

import (
	"tutorials/chat_server/protocol"
)

type ChatClient interface {
	Dial(address string) error
	Send(command interface{}) error
	SendMessage(message string) error
	SetName(name string) error
	Start()
	Close()
	Incoming() chan protocol.MessageCommand
}

