# Package waitfor

Waitfor waits while specified tcp ports will be available.

```
go get github.com/hummerd/gostuff/waitfor
```

Example for using WaitServcies:
``` go
import (
	"github.com/hummerd/gostuff/waitfor"
	"github.com/hummerd/gostuff/errors"
)

func WaitForServcies(mongoDBConn, postgreConn, mySqlconn, redisConn string) error {
	err := waitfor.WaitServcies(
		time.Minute, time.Millisecond*20, 
		mongoDBConn, 
		postgreConn, 
		mySqlconn,
		redisConn)
	return errors.Wrap(err, "Service not available: ")
}
```
