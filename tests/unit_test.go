package tests

import (
	"testing"
	"fmt"
	"github.com/straysh/go_mnemonic/utils"
	"strconv"
)

func Test_words(t *testing.T) {
	n := 280
	fmt.Println( utils.IntTo8BitsArray( n ) )
	fmt.Printf( "%10s\n", string(utils.IntTo8BitsString( n )) )
	fmt.Printf("%10s\n", strconv.FormatInt(280, 2) )
}
