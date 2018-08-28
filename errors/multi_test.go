package errors

import (
	"testing"
)

func TestMuiltiError(t *testing.T) {
	me := &MultiError{}
	me.Add(New("err1"))
	me.Add(nil)
	me.Add(New("err2"))
	me.Add(New("err3"))

	if me.ActualLen() != 3 {
		t.Fatal("Wrong error count")
	}

	if me.AddedLen() != 4 {
		t.Fatal("Wrong add count")
	}

	e := error(me)
	emsg := e.Error()
	if emsg != "err1; err2; err3; " {
		t.Fatal("Wrong err message")
	}
}

func TestNewMuiltiError(t *testing.T) {
	me := NewMultiError(New("err1"), New("err2"))

	if me.ActualLen() != 2 {
		t.Fatal("Wrong error count")
	}

	if me.AddedLen() != 2 {
		t.Fatal("Wrong add count")
	}

	e := error(me)
	emsg := e.Error()
	if emsg != "err1; err2; " {
		t.Fatal("Wrong err message")
	}
}
