package ioutil

import (
	"bytes"
	"io"
	"testing"
)

func TestCachedReaderLess(t *testing.T) {
	cr, err := NewPrefixReader(nil, 10)
	if err != nil {
		t.Fatal(err)
	}

	data := &bytes.Buffer{}
	s := "123456789012345678901234567890"
	data.WriteString(s)

	err = cr.Reset(data)
	if err != nil {
		t.Fatal(err)
	}

	p := cr.Prefix()
	if string(p) != s[:10] {
		t.Fatal("Wrong prefix ", p)
	}

	d := make([]byte, 12)

	read := 0
	n, err := cr.Read(d)
	if n != len(d) || string(d) != s[read:read+len(d)] || err != nil {
		t.Fatal("Can not read first ", n, err)
	}
	read += n

	n, err = cr.Read(d)
	if n != len(d) || string(d) != s[read:read+len(d)] || err != nil {
		t.Fatal("Can not read second ", n, err)
	}
	read += n

	n, err = cr.Read(d)
	if n != 6 || string(d[:n]) != s[read:read+6] || err != nil {
		t.Fatal("Can not read third ", n, err)
	}

	err = cr.Close()
	if err != nil {
		t.Fatal("Close err ", err)
	}
}

func TestCachedReader_WriterTo(t *testing.T) {
	cr, err := NewPrefixReader(nil, 10)
	if err != nil {
		t.Fatal(err)
	}

	var wt io.WriterTo = cr
	t.Log("WriterTo to implemented", wt)

	s := "123456789012345678901234567890"
	data := bytes.NewBufferString(s)
	err = cr.Reset(data)
	if err != nil {
		t.Fatal(err)
	}

	p := cr.Prefix()
	if string(p) != s[:10] {
		t.Fatal("Wrong prefix ", p)
	}

	buff := &bytes.Buffer{}
	cr.WriteTo(buff)

	if buff.String() != s {
		t.Fatal("Wrong data ", buff.String(), s)
	}

	// Case Read + WriterTo
	data = bytes.NewBufferString(s)
	err = cr.Reset(data)
	if err != nil {
		t.Fatal(err)
	}

	prefix := make([]byte, 2)
	n, err := cr.Read(prefix)
	if n != len(prefix) || err != nil {
		t.Fatal(err)
	}

	buff = &bytes.Buffer{}
	cr.WriteTo(buff)

	if buff.String() != s[len(prefix):] {
		t.Fatal("Wrong data ", buff.String(), s)
	}

	// Case Read more than prefix + WriterTo
	data = bytes.NewBufferString(s)
	err = cr.Reset(data)
	if err != nil {
		t.Fatal(err)
	}

	prefix = make([]byte, 22)
	n, err = cr.Read(prefix)
	if n != len(prefix) || err != nil {
		t.Fatal(err)
	}

	buff = &bytes.Buffer{}
	cr.WriteTo(buff)

	if buff.String() != s[len(prefix):] {
		t.Fatal("Wrong data ", buff.String(), s)
	}
}

func TestCachedReaderMore(t *testing.T) {
	cr, err := NewPrefixReader(nil, 100)
	if err != nil {
		t.Fatal(err)
	}

	data := &bytes.Buffer{}
	s := "123456789012345678901234567890"
	data.WriteString(s)

	_ = cr.Reset(data)

	p := cr.Prefix()
	if string(p) != s {
		t.Fatal("Wrong buffer ", p)
	}

	d := make([]byte, 12)

	read := 0
	n, err := cr.Read(d)
	if n != len(d) || string(d) != s[read:read+len(d)] || err != nil {
		t.Fatal("Can not read first ", n, err)
	}
	read += n

	n, err = cr.Read(d)
	if n != len(d) || string(d) != s[read:read+len(d)] || err != nil {
		t.Fatal("Can not read second ", n, err)
	}
	read += n

	n, err = cr.Read(d)
	if n != 6 || string(d[:n]) != s[read:read+6] || err != io.EOF {
		t.Fatal("Can not read third ", n, err)
	}
}

func TestCachedWriterLess(t *testing.T) {
	cw := NewPrefixWriter(nil, 10)

	data := &bytes.Buffer{}
	s := "123456789012345678901234567890"
	data.WriteString(s)

	buff := &bytes.Buffer{}
	cw.Reset(buff)

	wl := 12
	written := 0
	n, err := cw.Write([]byte(s[:wl]))
	if n != wl || err != nil {
		t.Fatal("Can not write first ", n, err)
	}
	written += n

	n, err = cw.Write([]byte(s[written : written+wl]))
	if n != wl || err != nil {
		t.Fatal("Can not write second ", n, err)
	}
	written += n

	n, err = cw.Write([]byte(s[written : written+6]))
	if n != 6 || err != nil {
		t.Fatal("Can not write third ", n, err)
	}
	written += n

	p := cw.Prefix()
	if string(p) != s[:10] {
		t.Fatal("Wrong prefix ", p)
	}

	if buff.String() != s {
		t.Fatal("Wrong data ", buff.String(), s)
	}

	err = cw.Close()
	if err != nil {
		t.Fatal("Close err ", err)
	}
}

func TestCachedWriterMore(t *testing.T) {
	cw := NewPrefixWriter(nil, 40)

	data := &bytes.Buffer{}
	s := "123456789012345678901234567890"
	data.WriteString(s)

	buff := &bytes.Buffer{}
	cw.Reset(buff)

	wl := 12
	written := 0
	n, err := cw.Write([]byte(s[:wl]))
	if n != wl || err != nil {
		t.Fatal("Can not write first ", n, err)
	}
	written += n

	n, err = cw.Write([]byte(s[written : written+wl]))
	if n != wl || err != nil {
		t.Fatal("Can not write second ", n, err)
	}
	written += n

	n, err = cw.Write([]byte(s[written : written+6]))
	if n != 6 || err != nil {
		t.Fatal("Can not write third ", n, err)
	}
	written += n

	p := cw.Prefix()
	if string(p) != s {
		t.Fatal("Wrong prefix ", p)
	}

	if buff.String() != s {
		t.Fatal("Wrong data ", buff.String(), s)
	}
}

func TestCachedWriter_ReaderFrom(t *testing.T) {
	cw := NewPrefixWriter(nil, 10)

	var rf io.ReaderFrom = cw
	t.Log("ReaderFrom to implemented", rf)

	buff := &bytes.Buffer{}
	cw.Reset(buff)

	s := "123456789012345678901234567890"
	data := bytes.NewBufferString(s)
	n, err := cw.ReadFrom(data)
	if n != int64(len(s)) || err != nil {
		t.Fatal("Can not write from reader ", n, err)
	}

	p := cw.Prefix()
	if string(p) != s[:10] {
		t.Fatal("Wrong prefix ", p)
	}

	if buff.String() != s {
		t.Fatal("Wrong data ", buff.String(), s)
	}

	// Write + ReadFrom
	buff = &bytes.Buffer{}
	cw.Reset(buff)

	nn, err := cw.Write([]byte(s[:2]))
	if nn != 2 || err != nil {
		t.Fatal("Failed to write", nn, err)
	}

	data = bytes.NewBufferString(s[2:])
	n, err = cw.ReadFrom(data)
	if n != int64(len(s)-2) || err != nil {
		t.Fatal("Can not write from reader ", n, err)
	}

	p = cw.Prefix()
	if string(p) != s[:10] {
		t.Fatal("Wrong prefix ", p)
	}

	if buff.String() != s {
		t.Fatal("Wrong data ", buff.String(), s)
	}

	// Write more than prefix + ReadFrom
	buff = &bytes.Buffer{}
	cw.Reset(buff)

	nn, err = cw.Write([]byte(s[:22]))
	if nn != 22 || err != nil {
		t.Fatal("Failed to write", nn, err)
	}

	data = bytes.NewBufferString(s[22:])
	n, err = cw.ReadFrom(data)
	if n != int64(len(s)-22) || err != nil {
		t.Fatal("Can not write from reader ", n, err)
	}

	p = cw.Prefix()
	if string(p) != s[:10] {
		t.Fatal("Wrong prefix ", p)
	}

	if buff.String() != s {
		t.Fatal("Wrong data ", buff.String(), s)
	}
}
