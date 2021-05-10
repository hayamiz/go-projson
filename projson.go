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
	stateArray0       // array with no member
	stateArray1       // array with more than one members
	stateObject0      // object with no member
	stateObject1      // object with more than one members
	stateObject0Keyed // object with key specified
	stateObject1Keyed // object with key specified
)

const (
	colorNormal  = 0
	colorBold    = 1
	colorBlack   = 30
	colorRed     = 31
	colorGreen   = 32
	colorYellow  = 33
	colorBlue    = 34
	colorMagenta = 35
	colorCyan    = 36
	colorWhite   = 37
)

const (
	colorKey    = colorRed
	colorInt    = colorGreen
	colorFloat  = colorCyan
	colorString = colorMagenta
)

type JsonPrinter struct {
	state     printerState
	pathStack *list.List
	buffer    *bytes.Buffer
	style     int
	termwid   int
	color     bool
	err       error

	// position in current line (used for smart style)
	linepos int
	curKey  string
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
		color:     false,
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
	printer.color = false
	printer.err = nil
	printer.linepos = 0
}

func (printer *JsonPrinter) Error() error {
	return printer.err
}

func (printer *JsonPrinter) SetStyle(style int) error {
	if printer.state != stateInit {
		return errors.New("Style cannot changed after putting some items")
	}

	printer.style = style
	return nil
}

func (printer *JsonPrinter) SetTermWidth(termwid int) error {
	if printer.state != stateInit {
		return errors.New("Terminal width cannot changed after putting some items")
	}

	printer.termwid = termwid
	return nil
}

func (printer *JsonPrinter) SetColor(color bool) error {
	if printer.state != stateInit {
		return errors.New("Color mode cannot changed after putting some items")
	}

	printer.color = color
	return nil
}

func (printer *JsonPrinter) String() (string, error) {
	if printer.state == stateInit || printer.state == stateFinal {
		return printer.buffer.String(), nil
	}

	return "", errors.New("Some object/array is not finished.")
}

func indent(str string, n int) string {
	buffer := bytes.NewBuffer([]byte{})
	for i := 0; i < n; i++ {
		buffer.WriteString(str)
	}

	return buffer.String()
}

func color(str string, colorcode int) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", colorcode, str)
}

func (printer *JsonPrinter) BeginArray() error {
	if printer.err != nil {
		return printer.err
	}

	switch printer.state {
	case stateInit: // OK
	case stateArray0: // OK
	case stateArray1: // OK
	case stateObject0Keyed: // OK
	case stateObject1Keyed: // OK
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

	var newchunk string
	var colorchunk string
	if printer.style == SmartStyle {
		if printer.state == stateObject0Keyed {
			newchunk = fmt.Sprintf("%s: [", printer.curKey)
			colorchunk = fmt.Sprintf("%s: [", color(printer.curKey, colorKey))
			if printer.color {
				printer.buffer.WriteString(colorchunk)
			} else {
				printer.buffer.WriteString(newchunk)
			}
			printer.linepos += len(newchunk)
			printer.curKey = ""
		} else if printer.state == stateObject1Keyed {
			newchunk = fmt.Sprintf(",\n%s%s: [", indent(" ", cur_level), printer.curKey)
			colorchunk = fmt.Sprintf(",\n%s%s: [", indent(" ", cur_level), color(printer.curKey, colorKey))
			if printer.color {
				printer.buffer.WriteString(colorchunk)
			} else {
				printer.buffer.WriteString(newchunk)
			}
			printer.linepos += len(newchunk) - 2 + len(indent(" ", cur_level))
			printer.curKey = ""
		} else if printer.state == stateInit || printer.state == stateArray0 {
			newchunk = fmt.Sprintf("[")
			printer.buffer.WriteString(newchunk)
			printer.linepos += len(newchunk)
		} else if printer.state == stateArray1 {
			newchunk = fmt.Sprintf(", [")
			printer.buffer.WriteString(newchunk)
			printer.linepos += len(newchunk)
		}

		if printer.linepos >= printer.termwid {
			printer.buffer.WriteString("\n" + indent(" ", cur_level+1))
			printer.linepos = cur_level + 1
		}
	} else {
		switch printer.state {
		case stateInit:
			printer.buffer.WriteString("[")
		case stateArray0:
			printer.buffer.WriteString("[")
		case stateArray1:
			printer.buffer.WriteString(",[")
		case stateObject0Keyed:
			if printer.color {
				printer.buffer.WriteString(fmt.Sprintf("%s:[", color(printer.curKey, colorKey)))
			} else {
				printer.buffer.WriteString(fmt.Sprintf("%s:[", printer.curKey))
			}
			printer.curKey = ""
		case stateObject1Keyed:
			if printer.color {
				printer.buffer.WriteString(fmt.Sprintf(",%s:[", color(printer.curKey, colorKey)))
			} else {
				printer.buffer.WriteString(fmt.Sprintf(",%s:[", printer.curKey))
			}
			printer.curKey = ""
		}
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
	printer.linepos += 1
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
		case int64:
			if err := printer.PutInt64(v.(int64)); err != nil {
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
	case stateObject0Keyed: // OK
	case stateObject1Keyed: // OK
	default:
		printer.err = errors.New("Cannot start object in this context")
		return printer.err
	}

	var cur_level int
	if printer.pathStack.Len() == 0 {
		cur_level = 0
	} else {
		cur_level = printer.pathStack.Back().Value.(*pathStackFrame).level
	}

	var newchunk string
	var colorchunk string
	if printer.style == SmartStyle {
		if printer.state == stateObject0Keyed {
			newchunk = fmt.Sprintf("%s: {", printer.curKey)
			colorchunk = fmt.Sprintf("%s: {", color(printer.curKey, colorKey))
			if printer.color {
				printer.buffer.WriteString(colorchunk)
			} else {
				printer.buffer.WriteString(newchunk)
			}
			printer.linepos = len(newchunk) - 1 + len(indent(" ", cur_level))
			printer.curKey = ""
		} else if printer.state == stateObject1Keyed {
			newchunk = fmt.Sprintf(",\n%s%s: {", indent(" ", cur_level), printer.curKey)
			colorchunk = fmt.Sprintf(",\n%s%s: {", indent(" ", cur_level), color(printer.curKey, colorKey))
			if printer.color {
				printer.buffer.WriteString(colorchunk)
			} else {
				printer.buffer.WriteString(newchunk)
			}
			printer.linepos = len(newchunk) - 2 + len(indent(" ", cur_level))
			printer.curKey = ""
		} else if printer.state == stateInit || printer.state == stateArray0 {
			newchunk = fmt.Sprintf("{")
			printer.buffer.WriteString(newchunk)
			printer.linepos += len(newchunk)
		} else if printer.state == stateArray1 {
			newchunk = fmt.Sprintf(", {")
			printer.buffer.WriteString(newchunk)
			printer.linepos += len(newchunk)
		}
	} else {
		switch printer.state {
		case stateInit:
			printer.buffer.WriteString("{")
		case stateArray0:
			printer.buffer.WriteString("{")
		case stateArray1:
			printer.buffer.WriteString(",{")
		case stateObject0Keyed:
			if printer.color {
				printer.buffer.WriteString(fmt.Sprintf("%s:{", color(printer.curKey, colorKey)))
			} else {
				printer.buffer.WriteString(fmt.Sprintf("%s:{", printer.curKey))
			}
			printer.curKey = ""
		case stateObject1Keyed:
			if printer.color {
				printer.buffer.WriteString(fmt.Sprintf(",%s:{", color(printer.curKey, colorKey)))
			} else {
				printer.buffer.WriteString(fmt.Sprintf(",%s:{", printer.curKey))
			}
			printer.curKey = ""
		}
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
	printer.linepos += 1
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

func (printer *JsonPrinter) putLiteral(literal string, colorliteral string) error {
	switch printer.state {
	case stateInit: // OK
	case stateArray0: // OK
	case stateArray1: // OK
	case stateObject0Keyed: // OK
	case stateObject1Keyed: // OK
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

	var newchunk string
	var colorchunk string
	commasep := false
	switch printer.state {
	case stateInit:
		newchunk = literal
		colorchunk = colorliteral
	case stateArray0:
		newchunk = literal
		colorchunk = colorliteral
	case stateArray1:
		commasep = true
		newchunk = literal
		colorchunk = colorliteral
	case stateObject0Keyed:
		if printer.style == SmartStyle {
			newchunk = fmt.Sprintf("%s: %s", printer.curKey, literal)
			colorchunk = fmt.Sprintf("%s: %s", color(printer.curKey, colorKey), colorliteral)
			printer.curKey = ""
		} else {
			newchunk = fmt.Sprintf("%s:%s", printer.curKey, literal)
			colorchunk = fmt.Sprintf("%s:%s", color(printer.curKey, colorKey), colorliteral)
			printer.curKey = ""
		}
	case stateObject1Keyed:
		commasep = true
		if printer.style == SmartStyle {
			newchunk = fmt.Sprintf("%s: %s", printer.curKey, literal)
			colorchunk = fmt.Sprintf("%s: %s", color(printer.curKey, colorKey), colorliteral)
			printer.curKey = ""
		} else {
			newchunk = fmt.Sprintf("%s:%s", printer.curKey, literal)
			colorchunk = fmt.Sprintf("%s:%s", color(printer.curKey, colorKey), colorliteral)
			printer.curKey = ""
		}
	}

	if printer.style == SmartStyle {
		commalen := 0
		if commasep {
			commalen = 2
		}

		if printer.linepos+len(newchunk)+commalen >= printer.termwid+1 {
			if commasep {
				printer.buffer.WriteString(",\n")
				printer.buffer.WriteString(indent(" ", cur_level))
				printer.linepos = cur_level
			}
		} else {
			if commasep {
				printer.buffer.WriteString(", ")
				printer.linepos += 2
			}
		}

		if printer.color {
			printer.buffer.WriteString(colorchunk)
		} else {
			printer.buffer.WriteString(newchunk)
		}

		printer.linepos += len(newchunk)
	} else {
		if commasep {
			printer.buffer.WriteString(",")
		}
		if printer.color {
			printer.buffer.WriteString(colorchunk)
		} else {
			printer.buffer.WriteString(newchunk)
		}
		printer.linepos += len(newchunk)
	}

	// state transitions
	switch printer.state {
	case stateInit:
		printer.state = stateFinal
	case stateArray0:
		printer.state = stateArray1
	case stateObject0Keyed:
		printer.state = stateObject1
	case stateObject1Keyed:
		printer.state = stateObject1
	}

	return nil
}

func (printer *JsonPrinter) PutInt(v int) error {
	str := strconv.Itoa(v)
	return printer.putLiteral(str, color(str, colorInt))
}

func (printer *JsonPrinter) PutInt64(v int64) error {
	str := strconv.FormatInt(v, 10)
	return printer.putLiteral(str, color(str, colorInt))
}

func (printer *JsonPrinter) PutFloat(v float64) error {
	str := strconv.FormatFloat(v, 'f', -1, 64)
	return printer.putLiteral(str, color(str, colorFloat))
}

func (printer *JsonPrinter) PutFloatFmt(v float64, fmtstr string) error {
	str := fmt.Sprintf(fmtstr, v)
	return printer.putLiteral(str, color(str, colorFloat))
}

func (printer *JsonPrinter) PutString(v string) error {
	vs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	str := string(vs)

	return printer.putLiteral(str, color(str, colorString))
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

	vss := string(vs)

	printer.curKey = vss

	if printer.state == stateObject0 {
		printer.state = stateObject0Keyed
	} else {
		printer.state = stateObject1Keyed
	}

	return nil
}
