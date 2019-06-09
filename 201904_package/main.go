package main

import "github.com/keiji0/sandbox/201904_package/b"

// mainからだとa/bなのでinternalは見えない
// import "github.com/keiji0/sandbox/package/a/b/internal"

func main() {
	b.Hoge()
	// internal.AAA()
}
