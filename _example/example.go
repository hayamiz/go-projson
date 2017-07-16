//usr/bin/env go run $0 $@ ; exit

package main

import (
	"fmt"
	"go-projson"
)

func main() {
	jp := projson.NewPrinter()

	jp.BeginArray()
	jp.PutInt(10)
	jp.PutInt(20)
	jp.PutString("hello")
	jp.PutString("double quote \" string")
	jp.BeginArray()
	jp.PutInt(30)
	jp.PutInt(40)
	jp.FinishArray()
	jp.PutInt(50)
	jp.FinishArray()
	fmt.Println(jp.String())
}
