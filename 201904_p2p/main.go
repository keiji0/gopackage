package main

import (
	"log"
	"net"

	"github.com/keiji0/sandbox/201904_p2p/net/message"
)

func main() {

	clientChan := make(chan string)
	serverChan := make(chan string)
	go clientStart(clientChan)
	go serverStart(serverChan)

	for message := range serverChan {
		log.Printf("message %s", message)
		switch message {
		case "start":
			clientChan <- message
		}
	}
}

func clientStart(ch <-chan string) {
	<-ch
	log.Println("client start")

	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	hello := message.NewHello()
	message.Send(conn, hello)

	log.Println("client end")
}

func serverStart(ch chan<- string) {
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	ch <- "start"
	conn, err := listener.Accept()
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	for {
		message, err := message.Receive(conn)
		log.Println("server receive")
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(message)
	}
}
