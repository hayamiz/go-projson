package projson

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	style     int
	termwid   int
	err       error

	// position in current line (used for smart style)
	linepos int
}

type frameType int

const (
	frameArray frameType = iota
	frameObject
)

const (
	SimpleStyle int = iota
	SmartStyle
	PrettyStyle
)

type pathStackFrame struct {
	typ   frameType
	level int
}

func getSystemTermWidth() int {
	var wid, hei int

	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()

	if err != nil {
		return 80
	}

	fmt.Sscanf(string(out), "%d %d", &hei, &wid)

	return wid
}

func NewPrinter() *JsonPrinter {
	printer := &JsonPrinter{
		state:     stateInit,
		pathStack: list.New(),
		buffer:    bytes.NewBuffer([]byte{}),
		style:     SimpleStyle,
		termwid:   getSystemTermWidth(),
		err:       nil,
		linepos:   0,
	}

	return printer
}

func (printer *JsonPrinter) Reset() {
	printer.state = stateInit
	printer.pathStack = list.New()
	printer.buffer = bytes.NewBuffer([]byte{})
	printer.style = SimpleStyle
	printer.termwid = getSystemTermWidth()
	printer.err = nil
	printer.linepos = 0
}

func (printer *JsonPrinter) Error() error {
	return printer.err
}

func (printer *JsonPrinter) SetStyle(style int, termwid int) error {
	if printer.state != stateInit {
		return errors.New("Style cannot changed after putting some items")
	}

	printer.style = style
	printer.termwid = termwid
	return nil
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

	var cur_level int
	if printer.pathStack.Len() == 0 {
		cur_level = 0
	} else {
		cur_level = printer.pathStack.Back().Value.(*pathStackFrame).level
	}

	if printer.state == stateArray1 {
		printer.buffer.WriteString(",")
		printer.linepos += 1
	}

	if printer.style == SmartStyle {
		if printer.state == stateArray1 || printer.state == stateObjectKeyed {
			if cur_level > 0 {
				printer.buffer.WriteString("\n")
			}
			printer.linepos = 0
			for i := 0; i < cur_level; i++ {
				printer.buffer.WriteString(" ")
				printer.linepos += 1
			}
		}
	}
	printer.buffer.WriteString("[")
	printer.linepos += 1

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

	cur_level := printer.pathStack.Back().Value.(*pathStackFrame).level

	if printer.style == SmartStyle {
		if cur_level == 1 {
			if printer.linepos+1 > printer.termwid {
				printer.buffer.WriteString("\n")
				printer.linepos = 0
				printer.buffer.WriteString(" ")
				printer.linepos += 1
			}
		} else {
			if printer.linepos+2 > printer.termwid {
				printer.buffer.WriteString("\n")
				printer.linepos = 0
				for i := 0; i < cur_level; i++ {
					printer.buffer.WriteString(" ")
					printer.linepos += 1
				}
			}
		}
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

func (printer *JsonPrinter) PutArray(arr []interface{}) error {
	if err := printer.BeginArray(); err != nil {
		return err
	}

	for _, v := range arr {
		switch v.(type) {
		case int:
			if err := printer.PutInt(v.(int)); err != nil {
				return err
			}
		case float64:
			if err := printer.PutFloat(v.(float64)); err != nil {
				return err
			}
		case string:
			if err := printer.PutString(v.(string)); err != nil {
				return err
			}
		default:
			return errors.New("unknown type in array")
		}
	}

	if err := printer.FinishArray(); err != nil {
		return err
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

func (printer *JsonPrinter) PutObject(m map[string]interface{}) error {
	if err := printer.BeginObject(); err != nil {
		return err
	}

	for k, v := range m {
		if err := printer.PutKey(k); err != nil {
			return err
		}

		switch v.(type) {
		case int:
			if err := printer.PutInt(v.(int)); err != nil {
				return err
			}
		case float64:
			if err := printer.PutFloat(v.(float64)); err != nil {
				return err
			}
		case string:
			if err := printer.PutString(v.(string)); err != nil {
				return err
			}
		default:
			return errors.New("unknown type in array")
		}
	}

	if err := printer.FinishObject(); err != nil {
		return err
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

	var cur_level int
	if printer.pathStack.Len() == 0 {
		cur_level = 0
	} else {
		cur_level = printer.pathStack.Back().Value.(*pathStackFrame).level
	}

	if printer.state == stateArray0 {
		printer.state = stateArray1
	} else if printer.state == stateArray1 {
		if printer.style == SmartStyle {
			if printer.linepos+2+len(literal) >= printer.termwid {
				printer.buffer.WriteString(",\n")
				printer.linepos = 0
				for i := 0; i < cur_level; i++ {
					printer.buffer.WriteString(" ")
					printer.linepos += 1
				}
			} else {
				printer.buffer.WriteString(", ")
				printer.linepos += 2
			}
		} else {
			printer.buffer.WriteString(",")
			printer.linepos += 1
		}
		printer.state = stateArray1
	} else if printer.state == stateInit {
		printer.state = stateFinal
	} else if printer.state == stateObjectKeyed {
		printer.state = stateObject1
		if printer.style == SmartStyle {
			if printer.linepos+1+len(literal) >= printer.termwid {
				printer.buffer.WriteString("\n")
				printer.linepos = 0
				for i := 0; i < cur_level; i++ {
					printer.buffer.WriteString(" ")
					printer.linepos += 1
				}
			} else {
				printer.buffer.WriteString(" ")
				printer.linepos += 1
			}
		}
	}

	printer.buffer.WriteString(literal)
	printer.linepos += len(literal)

	return nil
}

func (printer *JsonPrinter) PutInt(v int) error {
	return printer.putLiteral(strconv.Itoa(v))
}

func (printer *JsonPrinter) PutFloat(v float64) error {
	return printer.putLiteral(strconv.FormatFloat(v, 'f', -1, 64))
}

func (printer *JsonPrinter) PutFloatFmt(v float64, fmtstr string) error {
	return printer.putLiteral(fmt.Sprintf(fmtstr, v))
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

	cur_level := printer.pathStack.Back().Value.(*pathStackFrame).level

	vss := string(vs)

	if printer.state != stateObject0 {
		if printer.style == SmartStyle {
			if printer.linepos+3+len(vss) > printer.termwid {
				printer.buffer.WriteString(",\n")
				printer.linepos = 0
				for i := 0; i < cur_level; i++ {
					printer.buffer.WriteString(" ")
					printer.linepos += 1
				}
			}
		} else {
			printer.buffer.WriteString(",")
			printer.linepos += 1
		}
	}
	printer.buffer.WriteString(vss)
	printer.linepos += len(vss)
	printer.buffer.WriteString(":")
	printer.linepos += 1
	printer.state = stateObjectKeyed

	return nil
}
