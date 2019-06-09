package b

import (
	"fmt"

	"github.com/keiji0/sandbox/package/a/b/internal"
)

func B() {
	fmt.Println("B()")
	internal.AAA()
	// foo.AAA()
}
