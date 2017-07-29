//usr/bin/env go run $0 $@ ; exit

package projson

import "testing"

func TestInt(t *testing.T) {
	var err error
	jp := NewPrinter()

	err = jp.PutInt(12345)
	if err != nil {
		t.Error("expected: err == nil\nactual: err != nil")
		return
	}

	expected := "12345"
	actual, _ := jp.String()
	if actual != expected {
		t.Errorf("expected: %v\nactual: %v\n", expected, actual)
	}

	err = jp.PutInt(56789)
	if err == nil {
		t.Error("expected: err != nil\nactual: err == nil")
		return
	}
}

func TestFloat(t *testing.T) {
	var err error
	jp := NewPrinter()

	err = jp.PutFloat(0.5)
	if err != nil {
		t.Error("expected: err == nil\nactual: err != nil")
		return
	}

	expected := "0.5"
	actual, _ := jp.String()
	if actual != expected {
		t.Errorf("expected: %v\nactual: %v\n", expected, actual)
	}

	err = jp.PutFloat(1.234)
	if err == nil {
		t.Error("expected: err != nil\nactual: err == nil")
		return
	}
}

func TestString(t *testing.T) {
	var err error
	jp := NewPrinter()

	err = jp.PutString("hello, world")
	if err != nil {
		t.Error("err != nil")
	}

	expected := "\"hello, world\""
	actual, _ := jp.String()
	if actual != expected {
		t.Errorf("expected: %v\nactual: %v\n", expected, actual)
	}

	err = jp.PutInt(56789)
	if err == nil {
		t.Error("expected: err != nil\nactual: err == nil")
	}
}

func TestReset(t *testing.T) {
	jp := NewPrinter()

	jp.PutInt(42)
	str, err := jp.String()
	if str != "42" {
		t.Errorf("expected: 42, actual: %s", str)
	}

	jp.Reset()

	jp.PutInt(1192)
	str, err = jp.String()
	if err != nil {
		t.Errorf("expected: err == nil, actual: err != nil")
	}
	if str != "1192" {
		t.Errorf("expected: 1192, actual: %s", str)
	}
}

func expectNil(t *testing.T, v interface{}) {
	if v != nil {
		t.Errorf("expected: nil, actual: %v", v)
	}
}

func expectNonNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Errorf("expected: non-nil, actual: nil")
	}
}

func TestArraySimple(t *testing.T) {
	var err error
	jp := NewPrinter()

	err = jp.FinishArray()
	if err == nil {
		t.Errorf("expected: non-nil, actual: nil")
	}
	if jp.Error() == nil {
		t.Errorf("expected: non-nil, actual: nil")
	}

	jp.Reset()

	if jp.Error() != nil {
		t.Errorf("expected: nil, actual: non-nil")
	}

	jp.BeginArray()
	jp.PutInt(1)
	jp.PutString("two")
	jp.PutInt(3)

	_, err = jp.String()
	if err == nil {
		t.Errorf("error should be returned when calling String() for non-finished array")
	}

	jp.FinishArray()

	if jp.Error() != nil {
		t.Errorf("error should not happen.")
	}

	expected := `[1,"two",3]`
	actual, _ := jp.String()
	if expected != actual {
		t.Errorf("Unexpected JSON output\nexpected: %v\nactual: %v",
			expected, actual)
	}
}

func TestPutArraySimple(t *testing.T) {
	jp := NewPrinter()

	jp.PutArray([]interface{}{1, 2, 3, "foo", 4.5})

	expected := `[1,2,3,"foo",4.5]`
	if actual, _ := jp.String(); expected != actual {
		t.Errorf("Unexpected JSON output\nexpected: %v\nactual: %v",
			expected, actual)
	}

	jp.Reset()
	if err := jp.PutArray([]interface{}{complex128(1)}); err == nil {
		t.Errorf("Should not accept values expect of int, string, and float64")
	}
}

func TestArrayEmpty(t *testing.T) {
	jp := NewPrinter()

	jp.BeginArray()
	jp.FinishArray()

	if jp.Error() != nil {
		t.Errorf("JSON printing failed.")
	}

	expected := `[]`
	actual, _ := jp.String()
	if expected != actual {
		t.Errorf("Unexpected JSON output\nexpected: %v\nactual: %v",
			expected, actual)
	}
}

func TestArray(t *testing.T) {
	var err error
	jp := NewPrinter()

	err = jp.BeginArray()
	expectNil(t, err)

	err = jp.PutInt(1)
	expectNil(t, err)

	err = jp.PutInt(2)
	expectNil(t, err)

	err = jp.PutString("hello world")
	expectNil(t, err)

	err = jp.BeginArray()
	expectNil(t, err)

	err = jp.PutInt(4)
	expectNil(t, err)

	err = jp.BeginArray()
	expectNil(t, err)

	err = jp.PutInt(5)
	expectNil(t, err)

	err = jp.PutInt(6)
	expectNil(t, err)

	err = jp.FinishArray()
	expectNil(t, err)

	err = jp.PutString("fo\"o")
	expectNil(t, err)

	err = jp.FinishArray()
	expectNil(t, err)

	err = jp.FinishArray()
	expectNil(t, err)

	expected := `[1,2,"hello world",[4,[5,6],"fo\"o"]]`
	actual, _ := jp.String()

	if expected != actual {
		t.Errorf("expected: %v\nactual: %v\n", expected, actual)
	}
}

func TestArray2(t *testing.T) {
	jp := NewPrinter()
	jp.BeginArray()
	{
		jp.PutInt(1)
		jp.PutString("two")
		jp.PutFloat(3.5)
	}
	jp.FinishArray()
	expectNil(t, jp.Error())
	expected := `[1,"two",3.5]`
	actual, _ := jp.String()
	if expected != actual {
		t.Errorf("expected: %v\nactual: %v\n", expected, actual)
	}

	jp = NewPrinter()
	jp.BeginArray()
	{
		jp.PutString("two")
		jp.PutFloat(3.5)
		jp.PutInt(1)
	}
	jp.FinishArray()
	expectNil(t, jp.Error())
	expected = `["two",3.5,1]`
	actual, _ = jp.String()
	if expected != actual {
		t.Errorf("expected: %v\nactual: %v\n", expected, actual)
	}

	jp = NewPrinter()
	jp.BeginArray()
	{
		jp.PutFloat(3.5)
		jp.PutInt(1)
		jp.PutString("two")
	}
	jp.FinishArray()
	expectNil(t, jp.Error())
	expected = `[3.5,1,"two"]`
	actual, _ = jp.String()
	if expected != actual {
		t.Errorf("expected: %v\nactual: %v\n", expected, actual)
	}
}

func TestObjectSimple(t *testing.T) {
	jp := NewPrinter()

	jp.BeginObject()
	jp.PutKey("key1")
	jp.PutInt(10)

	str, err := jp.String()
	if err == nil {
		t.Errorf("error should be returned when calling String() for non-finished object")
	}
	if str != "" {
		t.Errorf("empty string should be returned on error")
	}

	jp.FinishObject()

	if jp.Error() != nil {
		t.Errorf("JSON printing failed.")
	}

	expected := `{"key1":10}`
	actual, _ := jp.String()

	if expected != actual {
		t.Errorf("expected: %v\nactual: %v", expected, actual)
	}
}

func TestPutObjectSimple(t *testing.T) {
	jp := NewPrinter()

	jp.PutObject(map[string]interface{}{"key1": 10})

	expected := `{"key1":10}`
	if actual, _ := jp.String(); expected != actual {
		t.Errorf("expected: %v\nactual: %v", expected, actual)
	}
}

func TestObjectEmpty(t *testing.T) {
	jp := NewPrinter()

	jp.BeginObject()
	jp.FinishObject()

	if jp.Error() != nil {
		t.Errorf("JSON printing failed.")
	}

	expected := `{}`
	actual, _ := jp.String()
	if expected != actual {
		t.Errorf("Unexpected JSON output\nexpected: %v\nactual: %v",
			expected, actual)
	}
}

func TestObject(t *testing.T) {
	jp := NewPrinter()

	jp.BeginObject()
	{
		jp.PutKey("key1")
		jp.BeginObject()
		{
			jp.PutKey("key2")
			jp.BeginObject()
			{
				jp.PutKey("key3")
				jp.BeginObject()
				{
					jp.PutKey("key4")
					jp.PutString("value4")
				}
				jp.FinishObject()
				jp.PutKey("key5")
				jp.PutInt(123)
			}
			jp.FinishObject()
		}
		jp.FinishObject()
	}
	jp.FinishObject()

	if jp.Error() != nil {
		t.Errorf("JSON printing failed.")
	}

	expected := `{"key1":{"key2":{"key3":{"key4":"value4"},"key5":123}}}`
	actual, _ := jp.String()

	if expected != actual {
		t.Errorf("expected: %v\nactual: %v", expected, actual)
	}
}

func TestSetStyle(t *testing.T) {
	jp := NewPrinter()

	err := jp.SetStyle(SimpleStyle)
	if err != nil {
		t.Error("expected: err == nil, actual: err != nil")
	}
	err = jp.SetTermWidth(80)
	if err != nil {
		t.Error("expected: err == nil, actual: err != nil")
	}

	err = jp.SetStyle(SmartStyle)
	if err != nil {
		t.Error("expected: err == nil, actual: err != nil")
	}
	err = jp.SetStyle(PrettyStyle)
	if err != nil {
		t.Error("expected: err == nil, actual: err != nil")
	}

	jp.PutInt(42)

	err = jp.SetStyle(SimpleStyle)
	if err == nil {
		t.Error("SetStyle should return error after putting items")
	}
}

func TestArraySimpleSmartStyle(t *testing.T) {
	jp := NewPrinter()

	jp.SetStyle(SmartStyle)
	jp.SetTermWidth(10)

	jp.BeginArray()
	jp.PutInt(10)
	jp.PutInt(20)
	jp.PutInt(30)
	jp.PutFloat(4.5)
	jp.PutInt(50)
	jp.PutFloat(60.5)
	jp.FinishArray()

	expected := `[10, 20,
 30, 4.5,
 50, 60.5]`
	actual, _ := jp.String()

	if expected != actual {
		t.Errorf("expected: %v\nactual: %v", expected, actual)
	}

	jp.Reset()
	jp.SetStyle(SmartStyle)
	jp.SetTermWidth(10)

	jp.BeginArray()
	jp.PutString("1234567890")
	jp.PutInt(10)
	jp.PutInt(20)
	jp.BeginArray()
	jp.PutInt(1)
	jp.PutInt(2)
	jp.PutInt(3)
	jp.PutInt(4)
	jp.PutInt(5)
	jp.PutInt(6)
	jp.FinishArray()
	jp.BeginArray()
	jp.PutInt(1)
	jp.PutInt(234)
	jp.FinishArray()
	jp.BeginArray()
	jp.BeginArray()
	jp.BeginArray()
	jp.PutInt(1)
	jp.FinishArray()
	jp.FinishArray()
	jp.FinishArray()
	jp.FinishArray()

	expected = `["1234567890",
 10, 20, [
  1, 2, 3,
  4, 5, 6
  ], [1,
  234], [[
   [1]]]]`
	actual, _ = jp.String()

	if expected != actual {
		t.Errorf("expected: %v\nactual: %v", expected, actual)
	}
}

func TestObjectSimpleSmartStyle(t *testing.T) {
	jp := NewPrinter()

	jp.SetStyle(SmartStyle)
	jp.SetTermWidth(10)

	jp.BeginObject()
	jp.PutKey("key")
	jp.PutString("val")
	jp.PutKey("k")
	jp.PutString("v")
	jp.PutKey("k")
	jp.BeginArray()
	jp.PutInt(1)
	jp.PutInt(2)
	jp.FinishArray()
	jp.FinishObject()

	expected := `{"key": "val",
 "k": "v",
 "k": [
  1, 2]}`
	actual, _ := jp.String()

	if expected != actual {
		t.Errorf("expected: %v\nactual: %v", expected, actual)
	}

	jp.Reset()
	jp.SetTermWidth(80)
	jp.SetStyle(SmartStyle)

	jp.BeginObject()
	jp.PutKey("key1")
	jp.PutString("val1")
	jp.PutKey("key2")
	jp.PutString("val2")
	jp.PutKey("key3")
	jp.PutArray([]interface{}{1, 2, 3, 4})
	jp.PutKey("key4")
	jp.BeginObject()
	jp.PutKey("key5")
	jp.PutInt(5)
	jp.FinishObject()
	jp.FinishObject()

	expected = `{"key1": "val1", "key2": "val2",
 "key3": [1, 2, 3, 4],
 "key4": {"key5": 5}}`
	actual, _ = jp.String()

	if expected != actual {
		t.Errorf("expected: %v\nactual: %v", expected, actual)
	}

}

func TestObjectSimpleSmartStyle2(t *testing.T) {
	jp := NewPrinter()

	jp.SetStyle(SmartStyle)
	jp.SetTermWidth(80)

	jp.BeginObject()
	jp.PutKey("key")
	jp.PutString("val")
	jp.PutKey("k")
	jp.PutString("v")
	jp.PutKey("k")
	jp.BeginArray()
	jp.PutInt(1)
	jp.PutInt(2)
	jp.FinishArray()
	jp.FinishObject()

	expected := `{"key": "val", "k": "v",
 "k": [1, 2]}`
	actual, _ := jp.String()

	if expected != actual {
		t.Errorf("expected: %v\nactual: %v", expected, actual)
	}

	jp.Reset()
	jp.SetStyle(SmartStyle)
	jp.SetTermWidth(80)

	jp.BeginObject()
	jp.PutKey("key1")
	jp.BeginObject()
	jp.PutKey("key2")
	jp.PutArray([]interface{}{"elem1", "elem2", "elem3"})
	jp.PutKey("key3")
	jp.PutObject(map[string]interface{}{"key4": 1.5, "key5": 1.5})
	jp.FinishObject()
	jp.FinishObject()

	expected = `{"key1": {"key2": ["elem1", "elem2", "elem3"],
  "key3": {"key4": 1.5, "key5": 1.5}}}`
	actual, _ = jp.String()

	if expected != actual {
		t.Errorf("expected: %v\nactual: %v", expected, actual)
	}
}

func TestPutFloatFmt(t *testing.T) {
	p := NewPrinter()

	p.PutFloatFmt(1.2345, "%.2f")

	expected := "1.23"
	if actual, _ := p.String(); actual != expected {
		t.Errorf("expected: %v\nactual: %v", expected, actual)
	}
}
