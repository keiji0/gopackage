package main

// Socketのライフサイクルの確認

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"time"
)

var message = make(chan string)

func countServer(addr *net.TCPAddr) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	message <- "start"
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		defer conn.Close()
		for no := 0; true; no++ {
			if _, err := fmt.Fprintf(conn, "%d", no); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second)

			if 0 < no && (no%4 == 0) {
				buf, err := ioutil.ReadAll(conn)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Printf("server %s\n", buf)
			}
		}
	}
}

func testServer(addr *net.TCPAddr) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	message <- "start"
	for {
		fmt.Println("server accept")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			conn.Close()
		}()
		for no := 0; ; no++ {
			if _, err := fmt.Fprintf(conn, "%d", no); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second * time.Duration(rand.Intn(6)))
		}
	}
}

// 途中で切れるとどうなる
func testServerStop(addr *net.TCPAddr) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	message <- "start"
	for {
		fmt.Println("server accept")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			conn.Close()
		}()
		for no := 0; ; no++ {
			if _, err := fmt.Fprintf(conn, "%d", no); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second * time.Duration(rand.Intn(6)))
			return
		}
	}
}

func main() {

	rand.Seed(42)

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
			fmt.Printf("%d, %s\n", n, string(buf[:n]))
		}
		{
			buf := make([]byte, 100)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("%d, %s\n", n, string(buf))
		}

		// go func() {
		// 	for {
		// 		buf := make([]byte, 2)
		// 		_, err := conn.Read(buf)
		// 		if err != nil {
		// 			log.Fatalln(err)
		// 		}
		// 		recv <- buf
		// 	}
		// }()

		// for {
		// 	select {
		// 	case buf := <-recv:
		// 		fmt.Printf("clinet read size [%s]\n", buf)
		// 	case <-time.After(time.Second * 5):
		// 		fmt.Println("end")
		// 		return
		// 	}
		// }

		// i := 0
		// for {
		// 	buf := make([]byte, 20)
		// 	_, err = conn.Read(buf)
		// 	if err != nil {
		// 		time.Sleep(time.Second * 5)
		// 		log.Fatalln(err)
		// 	}
		// 	conn.Write([]byte("a"))

		// 	fmt.Printf("client [%s]\n", buf)
		// 	i++
		// 	if i == 4 {
		// 		conn.Close()
		// 	}
		// }
	}
}
