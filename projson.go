package projson

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"strconv"
)

type printerState int

const (
	stateInit printerState = iota
	stateFinal
	stateArray0      // array with no member
	stateArray1      // array with more than one members
	stateObject0     // object with no member
	stateObject1     // object with more than one members
	stateObjectKeyed // object with key specified
)

type JsonPrinter struct {
	state     printerState
	pathStack *list.List
	buffer    *bytes.Buffer
	err       error
}

type frameType int

const (
	frameArray frameType = iota
	frameObject
)

type pathStackFrame struct {
	typ   frameType
	level int
}

func NewPrinter() *JsonPrinter {
	printer := &JsonPrinter{
		state:     stateInit,
		pathStack: list.New(),
		buffer:    bytes.NewBuffer([]byte{}),
		err:       nil,
	}

	return printer
}

func (printer *JsonPrinter) Reset() {
	printer.state = stateInit
	printer.pathStack = list.New()
	printer.buffer = bytes.NewBuffer([]byte{})
	printer.err = nil
}

func (printer *JsonPrinter) Error() error {
	return printer.err
}

func (printer *JsonPrinter) String() (string, error) {
	if printer.state == stateInit || printer.state == stateFinal {
		return printer.buffer.String(), nil
	}

	return "", errors.New("Some object/array is not finished.")
}

func (printer *JsonPrinter) BeginArray() error {
	if printer.err != nil {
		return printer.err
	}

	switch printer.state {
	case stateInit: // OK
	case stateArray0: // OK
	case stateArray1: // OK
	case stateObjectKeyed: // OK
	default:
		printer.err = errors.New("Cannot start array in this context")
		return printer.err
	}

	if printer.state == stateArray1 {
		printer.buffer.WriteString(",")
	}
	printer.buffer.WriteString("[")

	var cur_level int
	if printer.pathStack.Len() == 0 {
		cur_level = 0
	} else {
		cur_level = printer.pathStack.Back().Value.(*pathStackFrame).level
	}

	printer.pathStack.PushBack(&pathStackFrame{typ: frameArray, level: cur_level + 1})
	printer.state = stateArray0

	return nil
}

func (printer *JsonPrinter) FinishArray() error {
	if printer.err != nil {
		return printer.err
	}

	switch printer.state {
	case stateArray0: // OK
	case stateArray1: // OK
	default:
		printer.err = errors.New("Cannot finish array in this context")
		return printer.err
	}

	if printer.pathStack.Len() == 0 ||
		printer.pathStack.Back().Value.(*pathStackFrame).typ != frameArray {
		printer.err = errors.New("No array stack frame found")
		return printer.err
	}

	printer.buffer.WriteString("]")
	printer.pathStack.Remove(printer.pathStack.Back())

	if printer.pathStack.Len() == 0 {
		printer.state = stateInit
	} else {
		switch printer.pathStack.Back().Value.(*pathStackFrame).typ {
		case frameArray:
			printer.state = stateArray1
		case frameObject:
			printer.state = stateObject1
		default:
			printer.err = errors.New("Cannot happen this case")
			return printer.err
		}
	}

	return nil
}

func (printer *JsonPrinter) BeginObject() error {
	if printer.err != nil {
		return printer.err
	}

	switch printer.state {
	case stateInit: // OK
	case stateArray0: // OK
	case stateArray1: // OK
	case stateObjectKeyed: // OK
	default:
		printer.err = errors.New("Cannot start object in this context")
		return printer.err
	}

	if printer.state == stateArray1 {
		printer.buffer.WriteString(",")
	}
	printer.buffer.WriteString("{")

	var cur_level int
	if printer.pathStack.Len() == 0 {
		cur_level = 0
	} else {
		cur_level = printer.pathStack.Back().Value.(*pathStackFrame).level
	}

	printer.pathStack.PushBack(&pathStackFrame{typ: frameObject, level: cur_level + 1})
	printer.state = stateObject0

	return nil
}

func (printer *JsonPrinter) FinishObject() error {
	if printer.err != nil {
		return printer.err
	}

	switch printer.state {
	case stateObject0: // OK
	case stateObject1: // OK
	default:
		printer.err = errors.New("Cannot finish object in this context")
		return printer.err
	}

	if printer.pathStack.Len() == 0 ||
		printer.pathStack.Back().Value.(*pathStackFrame).typ != frameObject {
		printer.err = errors.New("No object stack frame found")
		return printer.err
	}

	printer.buffer.WriteString("}")
	printer.pathStack.Remove(printer.pathStack.Back())

	if printer.pathStack.Len() == 0 {
		printer.state = stateInit
	} else {
		switch printer.pathStack.Back().Value.(*pathStackFrame).typ {
		case frameArray:
			printer.state = stateArray1
		case frameObject:
			printer.state = stateObject1
		default:
			printer.err = errors.New("Cannot happen this case")
			return printer.err
		}
	}

	return nil
}

func (printer *JsonPrinter) putLiteral(literal string) error {
	switch printer.state {
	case stateInit: // OK
	case stateArray0: // OK
	case stateArray1: // OK
	case stateObjectKeyed: // OK
	default:
		printer.err = errors.New("Cannot put literal (" + literal + ") in this context")
		return printer.err
	}

	if printer.state == stateArray0 {
		printer.state = stateArray1
	} else if printer.state == stateArray1 {
		printer.buffer.WriteString(",")
		printer.state = stateArray1
	} else if printer.state == stateInit {
		printer.state = stateFinal
	} else if printer.state == stateObjectKeyed {
		printer.state = stateObject1
	}

	printer.buffer.WriteString(literal)

	return nil
}

func (printer *JsonPrinter) PutInt(v int) error {
	return printer.putLiteral(strconv.Itoa(v))
}

func (printer *JsonPrinter) PutFloat(v float64) error {
	return printer.putLiteral(strconv.FormatFloat(v, 'f', -1, 64))
}

func (printer *JsonPrinter) PutString(v string) error {
	vs, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return printer.putLiteral(string(vs))
}

func (printer *JsonPrinter) PutKey(v string) error {
	switch printer.state {
	case stateObject0: // OK
	case stateObject1: // OK
	default:
		printer.err = errors.New("Cannot put key in this context")
		return printer.err
	}

	vs, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if printer.state != stateObject0 {
		printer.buffer.WriteString(",")
	}
	printer.buffer.WriteString(string(vs))
	printer.buffer.WriteString(":")
	printer.state = stateObjectKeyed

	return nil
}
