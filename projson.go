package projson

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
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
		pathStack: list.New(),
		buffer:    bytes.NewBuffer([]byte{}),
		err:       nil,
	}

	return printer
}

func (printer *JsonPrinter) Reset() {
	printer.pathStack = list.New()
	printer.buffer = bytes.NewBuffer([]byte{})
	printer.err = nil
}

func (printer *JsonPrinter) Error() error {
	return printer.err
}

func (printer *JsonPrinter) String() string {
	return printer.buffer.String()
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
		printer.err = errors.New("Cannot finish array ini this context")
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
			printer.state = stateObject0
		default:
			printer.err = errors.New("Cannot happen this case")
			return printer.err
		}
	}

	return nil
}

func (printer *JsonPrinter) PutInt(v int) error {
	switch printer.state {
	case stateInit: // OK
	case stateArray0: // OK
	case stateArray1: // OK
	default:
		printer.err = errors.New("Cannot put int in this context")
		return printer.err
	}

	if printer.state == stateArray0 || printer.state == stateInit {
		printer.buffer.WriteString(fmt.Sprintf("%d", v))

		if printer.state == stateArray0 {
			printer.state = stateArray1
		}
	} else {
		printer.buffer.WriteString(fmt.Sprintf(",%d", v))
	}

	if printer.state == stateInit {
		printer.state = stateFinal
	}

	return nil
}

func (printer *JsonPrinter) PutString(v string) error {
	switch printer.state {
	case stateInit: // OK
	case stateArray0: // OK
	case stateArray1: // OK
	default:
		printer.err = errors.New("Cannot put string in this context")
		return printer.err
	}

	vs, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if printer.state == stateArray0 || printer.state == stateInit {
		printer.buffer.WriteString(string(vs))

		if printer.state == stateArray0 {
			printer.state = stateArray1
		}
	} else {
		printer.buffer.WriteString(",")
		printer.buffer.WriteString(string(vs))
	}

	if printer.state == stateInit {
		printer.state = stateFinal
	}

	return nil
}
