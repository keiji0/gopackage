package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:12000")
	if err != nil {
		log.Fatalln(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept()
		log.Println("connect")
		if err != nil {
			log.Fatalln(err)
		}

		conn.Write([]byte("HTTP/1.0 200 OK\r\nContent-Type: text/html\r\nContent-Length: 7\r\n\r\n HELLO\r\n"))

		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(buf))

		conn.Close()
	}
}
