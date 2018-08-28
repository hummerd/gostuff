package errors

import (
	"strings"
	"sync"
	"testing"
)

func TestSyncMuiltiError(t *testing.T) {
	me := &SyncMultiError{}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		me.Add(New("err1"))
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		me.Add(New("err2"))
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		me.Add(nil)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		me.Add(New("err3"))
		wg.Done()
	}()

	wg.Wait()

	if me.ActualLen() != 3 {
		t.Fatal("Wrong error count")
	}

	if me.AddedLen() != 4 {
		t.Fatal("Wrong add count")
	}

	e := error(me)
	emsg := e.Error()
	if !strings.Contains(emsg, "err1; ") ||
		!strings.Contains(emsg, "err2; ") ||
		!strings.Contains(emsg, "err3; ") {
		t.Fatal("Wrong err message")
	}
}
