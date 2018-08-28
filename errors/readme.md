# Package errors

Helpers for idiomatic (I hope so) go error processing.

```
go get github.com/hummerd/gostuff/errors
```

Example for using New Wrap and Join:
``` go
import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hummerd/gostuff/errors"
)

func PrintFile(path string) (err error) {
    if path == "" {
        return errors.New("path can not be empty")
    }

    f, err := os.Open(path)
    if err != nil {
        return errors.Wrapf(err, "Can not open file %s: ", path)
    }
    // Not loosing both: err that could happen before Close and Close's error itself.
    defer func() { err = errors.Join(err, f.Close()) }()

    data, err := ioutil.ReadAll(f)
    if err != nil {
        return errors.Wrapf(err, "Can not read file %s: ", path)
    }

    _, err = fmt.Println(string(data))
    // Wrap returns nil if err is nil
    return errors.Wrap(err, "Can not print file content: ")
}
```

Example for using MultiError:
``` go
import (
	"io"

	"github.com/hummerd/gostuff/errors"
)

func WriteText(w io.Writer) error {
	if w == nil {
		return errors.New("w can not be nil")
	}

	me := &errors.MultiError{}
	_, err := io.WriteString(w, "String1")
	me.Add(err)
	_, err = io.WriteString(w, "String2")
	me.Add(err)
	_, err = io.WriteString(w, "String3")
	me.Add(err)

	return me.IfHasErrors()
	// Or more explicit version 
	// if me.ActualLen() > 0 {
	//     return me
	// }
	// return nil
}
```

Example for using SyncMultiError:
``` go
import (
	"io"
	"sync"

	"github.com/hummerd/gostuff/errors"
)

func MakeAsync(w io.WriterAt) error {
	if w == nil {
		return errors.New("w can not be nil")
	}

	me := &errors.SyncMultiError{}
	wg := sync.WaitGroup{}

	wg.Add(1)
	data1 := []byte("msg1")
	go func() {
		defer wg.Done()
		_, err := w.WriteAt(data1, 0)
		me.Add(err)
	}()

	wg.Add(1)
	off := int64(len(data1))
	data2 := []byte("msg2")
	go func() {
		defer wg.Done()
		_, err := w.WriteAt(data2, off)
		me.Add(err)
	}()

	wg.Wait()
	return me.IfHasErrors()
}
```