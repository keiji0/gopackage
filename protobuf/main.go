package main

import (
	"fmt"

	"github.com/keiji0/sandbox/protobuf/protocol"
)

func main() {
	fmt.Println("hoge")
	person := protocol.Person{}
	person.Email = "hoge"
	fmt.Println(person)
}
