//usr/bin/env go run $0 $@ ; exit

package main

import (
	"fmt"
	"go-projson"
)

func main() {
	jp := projson.NewPrinter()

	jp.PutInt(42)

	str, _ := jp.String()
	fmt.Println(str)

	jp.Reset()

	jp.BeginArray()
	{
		jp.PutInt(10)
		jp.PutInt(20)
		jp.PutString("hello")
		jp.PutString("double quote \" string")
		jp.BeginArray()
		{
			jp.PutInt(30)
			jp.PutInt(40)
		}
		jp.FinishArray()
		jp.PutInt(50)
	}
	jp.FinishArray()

	str, _ = jp.String()
	fmt.Println(str) // => [10,20,"hello","double quote \" string",[30,40],50]

	jp.Reset()

	jp.PutArray([]interface{}{1, 2, 3, 4.5, 5, "hoge"})
	str, _ = jp.String()
	fmt.Println(str) // => [1,2,3,4.5,5,"hoge"]

	jp.Reset()

	jp.PutObject(map[string]interface{}{"key1": 1, "key2": "str2", "key3": 4.56})
	str, _ = jp.String()
	fmt.Println(str) // => {"key1":1,"key2":"str2","key3":4.56}

	jp.Reset()
	jp.SetStyle(projson.SmartStyle)
	jp.PutObject(map[string]interface{}{"time": 10.0, "key1": "val1"})
	str, _ = jp.String()
	fmt.Println(str)

	jp.Reset()
	jp.SetStyle(projson.SmartStyle)
	jp.SetColor(true)
	jp.BeginObject()
	jp.PutKey("key1")
	jp.PutInt(12345)
	jp.PutKey("key2")
	jp.PutFloat(1.23)
	jp.PutKey("key3")
	jp.PutString("hello, world")
	jp.PutKey("key4")
	jp.PutArray([]interface{}{1, 2.34, "567"})
	jp.FinishObject()
	str, _ = jp.String()
	fmt.Println(str)
}
