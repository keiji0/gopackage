package b

import (
	"fmt"

	"github.com/keiji0/sandbox/201904_package/b/c"
)

func Hoge() {
	fmt.Println("b.Hoge()")
	c.Hoge()
}
