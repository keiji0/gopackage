package c

import "fmt"
import "github.com/keiji0/sandbox/201904_package/b/internal"

func Hoge() {
	fmt.Println("b.c Hoge()")
	internal.Hoge()
}
