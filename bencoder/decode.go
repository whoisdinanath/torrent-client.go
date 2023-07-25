package bencoder

import (
	"bufio"
	"io"
)

// decode a bencode stream, it parses the stream r
// and returns generic bencode object representation

func Decode(reader io.Reader) (data interface{}, err error) {
	bufioReader, ok := reader.(*bufio.Reader)
	if !ok {
		bufioReader = newBufioReader(reader)
		defer bufioReaderPool.Put(bufioReader)
	}
	return decodeFromReader(bufioReader)
}
