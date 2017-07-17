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
	actual := jp.String()
	if actual != expected {
		t.Errorf("expected: %v\nactual: %v\n", expected, actual)
	}

	err = jp.PutInt(56789)
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
	actual := jp.String()
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
	jp.FinishArray()

	if jp.Error() != nil {
		t.Errorf("error should not happen.")
	}

	expected := `[1,"two",3]`
	actual := jp.String()
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
	actual := jp.String()

	if expected != actual {
		t.Errorf("expected: %v\nactual: %v\n", expected, actual)
	}
}
