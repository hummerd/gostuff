package waitfor

import (
	"context"
	"database/sql"
	"net"
	"regexp"
	"time"

	"github.com/hummerd/gostuff/errors"
)

var (
	regexURL      = regexp.MustCompile(`^\w*://.*@(?P<host>.+):(?P<port>\d+).*`)           // something like amqp://user@host:port/some
	regexCredURL  = regexp.MustCompile(`^\w*://(?P<host>.+):(?P<port>\d+).*`)              // something like amqp://host:port/some
	regexHost     = regexp.MustCompile(`^(?P<host>[^\(\)]+):(?P<port>\d+)`)                // something like host:port
	regexDsn      = regexp.MustCompile(`\((?P<host>.+):(?P<port>\d+)\)`)                   // something like user:password@tcp(localhost:5555)/dbname?tls=skip-verify
	regexHostPort = regexp.MustCompile(`host\=(?P<host>.+?)\s.*port\=(?P<port>\d+)`)       // something like host=localhost port=1234
	regexPortHost = regexp.MustCompile(`port\=(?P<port>\d+)\s.*host\=(?P<host>.+?)(\s|$)`) // something like port=1234 host=localhost

	knownRegexp = []*regexp.Regexp{
		regexURL,
		regexCredURL,
		regexHost,
		regexDsn,
		regexHostPort,
		regexPortHost,
	}
)

func WaitSQL(timeout, retryAfter time.Duration, db *sql.DB) error {
	start := time.Now()
	left := timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		err := db.PingContext(ctx)
		if err == nil {
			return nil
		}

		time.Sleep(retryAfter)

		left = timeout - time.Since(start)
		if left < 0 {
			return errors.Wrap(err, "DB is not available: ")
		}
	}
}

// WaitTCPPort wait while it can connect to specified tcp port
func WaitTCPPort(timeout, retryAfter time.Duration, host, port string) error {
	start := time.Now()
	left := timeout
	for {
		conn, err := net.DialTimeout("tcp", host+":"+port, left)
		if conn != nil {
			_ = conn.Close()
			return nil
		}

		time.Sleep(retryAfter)

		left = timeout - time.Since(start)
		if left < 0 {
			return errors.Wrapf(err, "Servcie %s:%s not available", host, port)
		}
	}
}

// WaitServices waits for all specified services to be available.
// Service can be specified in one of the following forms:
//	 scheme://user@host:port/some
//   scheme://host:port/some
//   host:port
//   user:password@network(host:port)/path?etc
//   host=host port=port
//   port=port host=host
func WaitServices(timeout, retryAfter time.Duration, services ...string) error {
	start := time.Now()
	left := timeout
	for _, s := range services {
		h, p := parseConnectionString(s)
		if h == "" || p == "" {
			return errors.New("Can not parse service connection string: " + s)
		}

		err := WaitTCPPort(left, retryAfter, h, p)
		if err != nil {
			return err
		}
		left = timeout - time.Since(start)
	}

	return nil
}

func parseConnectionString(str string) (host, port string) {
	for _, r := range knownRegexp {
		sub := r.FindStringSubmatch(str)
		if len(sub) >= 3 {
			for i, n := range r.SubexpNames() {
				if i > 0 {
					if n == "host" {
						host = sub[i]
					} else if n == "port" {
						port = sub[i]
					}
				}
			}
			return
		}
	}
	return
}
