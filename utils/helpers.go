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
