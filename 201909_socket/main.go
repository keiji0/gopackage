package main

import (
	"time"

	"github.com/keiji0/sandbox/201909_socket/wire"
)

func main() {
	serv1 := wire.NewServer("localhost:22222")
	serv2 := wire.NewServer("localhost:22223")

	go serv1.Listen()
	go serv2.Listen()

	time.Sleep(time.Second * 1)
	serv2.Connect("localhost:22222")
}
