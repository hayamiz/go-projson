//usr/bin/env go run $0 $@ ; exit

package main

import (
	"fmt"
	"go-projson"
)

func main() {
	jp := projson.NewPrinter()
	jp.SetStyle(projson.SmartStyle)
	jp.SetTermWidth(10)

	jp.BeginArray()
	jp.PutString("1234567890")
	jp.FinishArray()

	str, _ := jp.String()
	fmt.Println(str)
}
