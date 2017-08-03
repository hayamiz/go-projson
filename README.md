# go-projson

[<img src="https://travis-ci.org/hayamiz/go-projson.svg?branch=master" />](https://travis-ci.org/hayamiz/go-projson)

Procedural JSON printer for Go

# Installation

```
$ go get github.com/hayamiz/go-projson
```

# Basic Usage

Basic usage of `go-projson` is:

1. Create `JsonPrinter` object by calling `projson.NewPrinter()` func.
2. Put JSON elements (int, float, string, object, array) one by one with following APIs:
  - `PutInt`, `PutFloat`, `PutString` ... functions for putting JSON primitive data.
  - `BeginArray`, `FinishArray` ... functions for putting arrays. Elements of an array are constructed by projson API calls between corresponding `BeginArray` and `FinishArray`.
  - `BeginObject`, `FinishObject` ... functions for putting objects. Members of an object are constructed by projson API calls between corresponding `BeginObject` and `FinishObject`, and each member must be keyed by a preceding `PutKey` API call.

## Example 1: basic usage

```
package main

import (
    projson "github.com/hayamiz/go-projson"
)

func main() {
    printer := projson.NewPrinter()

    // Building JSON output by calling Put*/Begin*/Finish* API functions
    printer.BeginObject()

    printer.PutKey("key1")
    printer.PutKey(12345)  // => "key1":12345

    printer.PutKey("key2")
    printer.BeginArray()
    printer.PutInt(12)
    printer.PutFloat(345.67)
    printer.PutString("hello, go-projson")
    printer.FinishArray()  // => "key2":[12, 345.67, "hello, go-projson"]

    printer.FinishObject()

    if str, err := printer.String(); err != nil {
        panic(err)
    } else {
        fmt.Println(str) // prints {"key1":12345,"key2":"key2":[12,345.67,"hello, go-projson"]}
    }
}
```

## Example 2: just print int, float, or string

```
    printer := projson.NewPrinter()
    printer.PutInt(42)
    str, _ := printer.String() // => 42
```

```
    printer := projson.NewPrinter()
    printer.PutFloat(123.45)

    str, _ := printer.String() // => 123.45
```

```
    printer := projson.NewPrinter()
    printer.PutString("hello, projson")

    str, _ := printer.String() // => "hello, projson"
```

## Example 2: nested array and nested object

```
    printer := projson.NewPrinter()

    printer.BeginArray()
      printer.PutInt(123)
      printer.BeginArray()
        printer.PutFloat(456.7)
        printer.BeginArray()
          printer.PutString("nest, nest, and nest")
        printer.FinishArray()
        printer.PutString("nest depth = two here")
      printer.FinishArray()
    printer.FinishArray()

    str, _ := printer.String() // => [123,[456.7,["nest, nest and nest"],"nest depth = two here"]]
```

```
    printer := projson.NewPrinter()

    printer.BeginObject()
      printer.PutKey("key1")
      printer.BeginObject()
        printer.PutKey("nested key1")
        printer.PutInt(12345)

        printer.PutKey("nested key2")
        printer.BeginObject()
          printer.PutKey("double nested key1")
          printer.PutFloat(678.9)
        printer.FinishObject()
      printer.FinishObject()
    printer.FinishObject()

    str, _ := printer.String() // => {"key1":{"nested key1":12345,"nested key2":{"double nested key1":678.9}}}
```


# License

MIT license

# Author

Yuto Hayamizu (hayamiz)
