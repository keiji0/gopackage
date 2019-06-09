package main

import (
	"fmt"
	"reflect"
)

func main() {
	type MyInt int

	myint := MyInt(8)

	var i interface{}
	i = myint
	fmt.Printf("%v\n", reflect.New(i))
}
