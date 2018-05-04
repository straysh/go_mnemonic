package tests

import (
	"testing"
	"fmt"
	"strconv"
	"github.com/straysh/go_mnemonic/utils"
)

func BenchmarkSprintfToBinary(b *testing.B) {
	for i:=0;i<b.N;i++ {
		fmt.Sprintf("%08b", 101)
	}
}

func BenchmarkFormatIntToBinary(b *testing.B) {
	for i:=0;i<b.N;i++ {
		strconv.FormatInt(int64(101), 2)
	}
}

func BenchmarkIntToBinary(b *testing.B) {
	for i:=0;i<b.N;i++ {
		utils.IntTo8BitsString( 101 )
	}
}
