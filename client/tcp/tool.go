package tcp

import (
	"bytes"
	"time"
)

func now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func WrapCodeString(code byte, msg string) *bytes.Buffer {
	var buff bytes.Buffer
	buff.Write([]byte{code, '#'})
	buff.Write([]byte(msg))
	return &buff
}

func WrapCodeByte(code byte, msg []byte) *bytes.Buffer {
	var buff bytes.Buffer
	buff.Write([]byte{code, '#'})
	buff.Write(msg)
	return &buff
}

func WrapCode(code byte) []byte {
	return []byte{code, '#'}
}

