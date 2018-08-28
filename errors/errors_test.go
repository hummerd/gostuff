package errors

import (
	"testing"
)

func TestError(t *testing.T) {
	e := New("err1")
	if e.Error() != "err1" {
		t.Fatal("Wrong error message", e.Error())
	}

	e = Newf("err1 %s", "test")
	if e.Error() != "err1 test" {
		t.Fatal("Wrong error message", e.Error())
	}

	e = Wrap(nil, "test")
	if e != nil {
		t.Fatal("Wrong nil wrap")
	}

	e = Wrap(New("err1"), "test ")
	if e.Error() != "test err1" {
		t.Fatal("Wrong wrap")
	}

	e = Wrapf(nil, "test %s", "msg")
	if e != nil {
		t.Fatal("Wrong nil wrapf")
	}

	e = Wrapf(New("err1"), "test %s ", "msg")
	if e.Error() != "test msg err1" {
		t.Fatal("Wrong wrap")
	}
}

func TestJointError(t *testing.T) {
	je := &JointError{}
	if je.Error() != "" {
		t.Fatal("Wrong empty joint error")
	}

	e1 := New("err1")
	e := Join(e1, nil)
	if e1 != e {
		t.Fatal("Wrong joint error", e.Error())
	}

	e2 := New("err2")
	e = Join(nil, e2)
	if e != e2 {
		t.Fatal("Wrong joint error", e.Error())
	}

	e = Join(e1, e2)
	je = e.(*JointError)
	if e1 != je.ErrOne || e2 != je.ErrTwo {
		t.Fatal("Wrong joint error", je.Error())
	}

	if e.Error() != "err1; also: err2" {
		t.Fatal("Wrong joint error", je.Error())
	}
}
