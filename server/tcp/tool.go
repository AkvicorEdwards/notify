package tcp

import (
	"bufio"
	"bytes"
	"notify/encryption"
	"os"
	"strings"
	"time"
)

func now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func WrapCodeString(code byte, msg string) []byte {
	var buff bytes.Buffer
	buff.Write([]byte{code, '#'})
	buff.Write(encryption.StringToBase64Byte(msg))
	return buff.Bytes()
}

func WrapCodeDoubleString(code byte, msg1, msg2 string) []byte {
	var buff bytes.Buffer
	buff.Write([]byte{code, '#'})
	buff.Write(encryption.StringToBase64Byte(msg1))
	buff.Write([]byte{'#'})
	buff.Write(encryption.StringToBase64Byte(msg2))
	return buff.Bytes()
}

func WrapDoubleCodeString(code1, code2 byte, msg string) []byte {
	var buff bytes.Buffer
	buff.Write([]byte{code1, '#'})
	buff.Write([]byte{code2, '#'})
	buff.Write(encryption.StringToBase64Byte(msg))
	return buff.Bytes()
}

func WrapTripleCodeString(code1, code2, code3 byte, msg string) []byte {
	var buff bytes.Buffer
	buff.Write([]byte{code1, '#'})
	buff.Write([]byte{code2, '#'})
	buff.Write([]byte{code3, '#'})
	buff.Write(encryption.StringToBase64Byte(msg))
	return buff.Bytes()
}
func WrapTripleCodeByte(code1, code2, code3 byte, msg []byte) []byte {
	var buff bytes.Buffer
	buff.Write([]byte{code1, '#'})
	buff.Write([]byte{code2, '#'})
	buff.Write([]byte{code3, '#'})
	buff.Write(encryption.ByteToBase64Byte(msg))
	return buff.Bytes()
}

func WrapCodeByte(code byte, msg []byte) []byte {
	var buff bytes.Buffer
	buff.Write([]byte{code, '#'})
	buff.Write(encryption.ByteToBase64Byte(msg))
	return buff.Bytes()
}

func WrapDoubleCodeByte(code1, code2 byte, msg []byte) []byte {
	var buff bytes.Buffer
	buff.Write([]byte{code1, '#'})
	buff.Write([]byte{code2, '#'})
	buff.Write(encryption.ByteToBase64Byte(msg))
	return buff.Bytes()
}

func WrapCode(code byte) []byte {
	return []byte{code, '#'}
}

func WrapDoubleCode(code1, code2 byte) []byte {
	return []byte{code1, '#', code2, '#'}
}

func GetTerminalInput() string {
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(input)
}
