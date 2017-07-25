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

func TestArraySimpleCompactStyle(t *testing.T) {
	jp := NewPrinter()

	jp.BeginArray()
	jp.FinishArray()
}
