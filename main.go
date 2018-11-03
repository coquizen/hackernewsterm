package main

import (
	"github.com/caninodev/hackernewsterm/client"
	"github.com/caninodev/hackernewsterm/server"
)

func main() {
	go server.CreateServer()
	go client.CreateClient()
}
