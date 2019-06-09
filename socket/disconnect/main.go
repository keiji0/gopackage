package main

// Socketのライフサイクルの確認

import (
	"fmt"
	"log"
	"net"
	"time"
)

var message = make(chan string)

// 途中で切れるとどうなる
func testServerStop(addr *net.TCPAddr) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	message <- "start"
	for {
		fmt.Println("server accept")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			conn.Close()
			fmt.Println("server close")
		}()
		for no := 0; ; no++ {
			if _, err := fmt.Fprintf(conn, "%d", no); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second * 1)
			return
		}
	}
}

func main() {

	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:12000")
	if err != nil {
		log.Fatalln(err)
	}

	// Server
	go func() {
		testServerStop(addr)
	}()

	// Client
	{
		<-message

		fmt.Println("client connect")
		conn, err := net.DialTCP("tcp4", nil, addr)
		if err != nil {
			log.Fatalln(err)
		}
		defer conn.Close()

		{
			buf := make([]byte, 100)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("read1 %d, %s\n", n, string(buf[:n]))
		}
		{
			buf := make([]byte, 100)
			n, err := conn.Read(buf)
			if err != nil {
				// closeされるとEOFが帰ってくる
				log.Fatalln(err)
			}
			fmt.Printf("read2 %d, %s\n", n, string(buf))
		}
	}
}
