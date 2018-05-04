package utils

import (
	"encoding/hex"
	"fmt"
)

func HexaAsBuffer(s string) []byte {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		panic(fmt.Errorf("hex.DecodeString fail@{%v}", s))
	}

	return decoded
}

func BufferToHexa(src []byte) string {
	encodedStr := hex.EncodeToString(src)
	return encodedStr
}

func IntTo8BitsArray(n int) [8]byte {
	var zero = byte('0')
	var one  = byte('1')
	var buffer [8]byte
	var a byte
	for i:= uint(1);i<=8;i++ {
		m := (n >> (i-1)) & 1
		switch m {
		case 0:
			buffer[8-i] = zero
		default:
			buffer[8-i] = one
		}
		a <<= 1
	}
	return buffer
}

func IntTo11BitsArray(n int) [11]byte {
	var zero = byte('0')
	var one  = byte('1')
	var buffer [11]byte
	var a byte
	for i:= uint(1);i<=11;i++ {
		m := (n >> (i-1)) & 1
		switch m {
		case 0:
			buffer[11-i] = zero
		default:
			buffer[11-i] = one
		}
		a <<= 1
	}
	return buffer
}

func IntTo8BitsString(n int) string {
	a := IntTo8BitsArray(n)
	return string(a[:])
}