package main

import "tutorials/chat_server/client/cmd"

func main() {
	var c cmd.ChatClient
	c = cmd.NewClient()
	c.Dial("localhost")

}
