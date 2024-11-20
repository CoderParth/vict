package main

import (
	"bytes"
	"compress/zlib"
)

func compress(data string) []byte {
	var buffer bytes.Buffer
	w := zlib.NewWriter(&buffer)
	w.Write([]byte(data))
	w.Close()
	return buffer.Bytes()
}
