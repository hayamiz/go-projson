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
		return
	}

	expected := "\"hello, world\""
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
