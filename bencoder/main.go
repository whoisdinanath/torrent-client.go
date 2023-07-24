package bencoder

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type builder interface {
	Int64(i int64)
	Uint64(i uint64)
	Float64(f float64)
	String(s string)
	Array()
	Map()

	Elem(i int) builder
	Key(s string) builder

	Flush()
}

type Reader interface {
	io.Reader
	io.ByteScanner
}

func decodeInt64(r *bufio.Reader, delim byte) (data int64, err error) {
	buf, err := readSlice(r, delim)
	if err != nil {
		return
	}
	data, err = strconv.ParseInt(string(buf), 10, 64)
	return
}

func readSlice(r *bufio.Reader, delim byte) (data []byte, err error) {
	if data, err = r.ReadSlice(delim); err != nil {
		return
	}
	lenData := len(data)
	if lenData > 0 {
		data = data[:lenData-1]
	} else {
		panic("bad r.ReadSlice() length")
	}
	return
}

func decodeString(r *bufio.Reader) (data string, err error) {
	length, err := decodeInt64(r, ':')
	if err != nil {
		return
	}
	if length < 0 {
		err = errors.New("Bad string length")
		return
	}
	if peekBuf, peekErr := r.Peek(int(length)); peekErr == nil {
		data = string(peekBuf)
		_, err = r.Discard(int(length))
		return
	}
	var buf = make([]byte, length)
	_, err = readFull(r, buf)
	if err != nil {
		return
	}
	data = string(buf)
	return

}

func readFull(r *bufio.Reader, buf []byte) (n int, err error) {
	fmt.Println("PASS")
	return
}
