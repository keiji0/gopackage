package main

// 巨大なメッセージをWriteするとどうなる?

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

var message = make(chan string)

// 途中で切れるとどうなる
func testServer(addr *net.TCPAddr) {
	rand.Seed(42)

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
		// buf := make([]byte, 1024*8*2)
		// for i := 0; i < cap(buf); i++ {
		// 	buf[i] = 1
		// }
		// n, err := conn.Write(buf)
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// fmt.Printf("sever: %d 書き込みました\n", n)
		var i int
		for i = 0; i < 1024*100; i++ {
			buf := []byte{1}
			_, err := conn.Write(buf)
			if err != nil {
				log.Fatalln(err)
			}
			// time.Sleep(time.Millisecond * time.Duration(rand.Intn(8)) * 10)
		}
		fmt.Printf("sever: %d 書き込みました\n", i)
	}
}

func main() {

	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:12000")
	if err != nil {
		log.Fatalln(err)
	}

	// Server
	go func() {
		testServer(addr)
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

		// buf := make([]byte, 9000)
		// n, err := conn.Read(buf)
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// fmt.Printf("client read size %d\n", n)

		time.Sleep(time.Second * 2)

		var i int
		for i = 0; i < 1024*8; i++ {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("client read %d\n", n)
		}

		time.Sleep(time.Second * 5)
	}
	fmt.Println("end")
}
