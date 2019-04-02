package main

import (
	"tutorials/chat_server/server/cmd"
)

func main() {
	var s cmd.ChatServer
	s = cmd.NewServer()
	s.Listen(":3333")

	s.Start()
}
