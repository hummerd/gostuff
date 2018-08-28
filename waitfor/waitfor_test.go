package waitfor

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestParseConnectionString(t *testing.T) {
	host, port := parseConnectionString("amqp://user:password@mydomain.com:5467/some")
	if host != "mydomain.com" || port != "5467" {
		t.Fatal("Wrong host port", host, port)
	}

	host, port = parseConnectionString("amqp://user@mydomain.com:5467/some")
	if host != "mydomain.com" || port != "5467" {
		t.Fatal("Wrong host port", host, port)
	}

	host, port = parseConnectionString("amqp://mydomain.com:5467/some")
	if host != "mydomain.com" || port != "5467" {
		t.Fatal("Wrong host port", host, port)
	}

	host, port = parseConnectionString("mydomain.com:5467")
	if host != "mydomain.com" || port != "5467" {
		t.Fatal("Wrong host port", host, port)
	}

	host, port = parseConnectionString("user:password@tcp(mydomain.com:5467)/dbname?tls=skip-verify")
	if host != "mydomain.com" || port != "5467" {
		t.Fatal("Wrong host port", host, port)
	}

	host, port = parseConnectionString("user:password@tcp(mydomain.com:5467)")
	if host != "mydomain.com" || port != "5467" {
		t.Fatal("Wrong host port", host, port)
	}

	host, port = parseConnectionString("(mydomain.com:5467)")
	if host != "mydomain.com" || port != "5467" {
		t.Fatal("Wrong host port", host, port)
	}

	host, port = parseConnectionString("some=wqer43 host=mydomain.com prm=456 port=5467")
	if host != "mydomain.com" || port != "5467" {
		t.Fatal("Wrong host port", host, port)
	}

	host, port = parseConnectionString("some=wqer43 host=mydomain.com prm=456 port=5467 arg=a")
	if host != "mydomain.com" || port != "5467" {
		t.Fatal("Wrong host port", host, port)
	}

	host, port = parseConnectionString("some=wqer43 port=5467  host=mydomain.com prm=456")
	if host != "mydomain.com" || port != "5467" {
		t.Fatal("Wrong host port", host, port)
	}

	host, port = parseConnectionString("some=wqer43 port=5467  host=mydomain.com")
	if host != "mydomain.com" || port != "5467" {
		t.Fatal("Wrong host port", host, port)
	}
}

func TestWait(t *testing.T) {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{net.IPv4(127, 0, 0, 1), 0, ""})
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	addrParts := strings.Split(l.Addr().String(), ":")
	err = WaitTCPPort(time.Second, time.Second, addrParts[0], addrParts[1])
	if err != nil {
		t.Fatal("Port not available", err)
	}
}

func TestWaitDelay(t *testing.T) {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{net.IPv4(127, 0, 0, 1), 0, ""})
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr()
	addrParts := strings.Split(addr.String(), ":")
	l.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		time.Sleep(time.Millisecond * 100)

		fmt.Println("Start listen", time.Now())
		tdpAddr, _ := net.ResolveTCPAddr(addr.Network(), addr.String())
		l, err := net.ListenTCP("tcp", tdpAddr)
		if err != nil {
			panic("Can not listen port")
		}
		wg.Wait()
		fmt.Println("Close", time.Now())
		defer l.Close()
	}()

	fmt.Println("WaitTCPPort", time.Now())
	err = WaitTCPPort(time.Second, time.Millisecond*100, addrParts[0], addrParts[1])
	wg.Done()
	if err != nil {
		t.Fatal("Port not available", err)
	}
}

func TestWaitFail(t *testing.T) {
	err := WaitTCPPort(time.Second, time.Millisecond*100, "localhost", "364589")
	if err == nil {
		t.Fatal("Fake port available")
	}
}

func TestWaitServices(t *testing.T) {
	lone, err := net.ListenTCP("tcp", &net.TCPAddr{net.IPv4(127, 0, 0, 1), 0, ""})
	if err != nil {
		t.Fatal(err)
	}
	addrPartsOne := strings.Split(lone.Addr().String(), ":")
	defer lone.Close()

	ltwo, err := net.ListenTCP("tcp", &net.TCPAddr{net.IPv4(127, 0, 0, 1), 0, ""})
	if err != nil {
		t.Fatal(err)
	}
	addrPartsTwo := strings.Split(ltwo.Addr().String(), ":")
	defer ltwo.Close()

	err = WaitServcies(
		time.Second, time.Millisecond*100,
		"localhost:"+addrPartsOne[1],
		"host=localhost port="+addrPartsTwo[1])
	if err != nil {
		t.Fatal("Services not available", err)
	}
}
