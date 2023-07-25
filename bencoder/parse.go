package bencoder

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
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
		err = errors.New("bad string length")
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
	return readAtLeast(r, buf, len(buf))
}

func readAtLeast(r *bufio.Reader, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, io.ErrShortBuffer
	}
	for n < min && err == nil {
		var nn int
		nn, err = r.Read(buf[n:])
		n += nn
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

func parseFromReader(r *bufio.Reader, build builder) (err error) {
	c, err := r.ReadByte()
	if err != nil {
		goto exit
	}
	switch {
	case c >= '0' && c <= '9':
		// this is a string
		err = r.UnreadByte()
		if err != nil {
			goto exit
		}
		var str string
		str, err = decodeString(r)
		if err != nil {
			goto exit
		}
		build.String(str)

	case c == 'd':
		// this is dictionary
		build.Map()
		for {
			c, err = r.ReadByte()
			if err != nil {
				goto exit
			}
			if c == 'e' {
				break
			}
			err = r.UnreadByte()
			if err != nil {
				goto exit
			}
			var key string
			key, err = decodeString(r)
			if err != nil {
				goto exit
			}
			// TODO: in pendantic mode, check for keys in ascending order.
			err = parseFromReader(r, build.Key(key))
			if err != nil {
				goto exit
			}
		}
	case c == 'i':
		var buf []byte
		buf, err = readSlice(r, 'e')
		if err != nil {
			goto exit
		}
		var str string
		var i int64
		var i2 uint64
		var f float64
		str = string(buf)
		if i, err = strconv.ParseInt(str, 10, 64); err == nil {
			build.Int64(i)
		} else if i2, err = strconv.ParseUint(str, 10, 64); err == nil {
			build.Uint64(i2)
		} else if f, err = strconv.ParseFloat(str, 64); err == nil {
			build.Float64(f)
		} else {
			err = errors.New("bad integer")
		}

	case c == 'l':
		// array
		build.Array()
		n := 0
		for {
			c, err = r.ReadByte()
			if err != nil {
				goto exit
			}
			if c == 'e' {
				break
			}
			err = r.UnreadByte()
			if err != nil {
				goto exit
			}
			err = parseFromReader(r, build.Elem(n))
			if err != nil {
				goto exit
			}
			n++
		}
	default:
		err = fmt.Errorf("unexpected character: '%v'", c)
	}
exit:
	build.Flush()
	return
}

// Parse parses the bencode stream and makes calls to
// the builder to construct a parsed representation.
func parse(reader io.Reader, builder builder) (err error) {
	//Check to see if the reader already fulfills the bufio.Reader interface.
	// Wrap it in a bufio.Reader if it doesn't.
	r, ok := reader.(*bufio.Reader)
	if !ok {
		r = newBufioReader(reader)
		defer bufioReaderPool.Put(r)
	}

	return parseFromReader(r, builder)

}

var bufioReaderPool sync.Pool

func newBufioReader(r io.Reader) *bufio.Reader {
	if v := bufioReaderPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReader(r)
}
