package tests

import (
	"testing"
	. "github.com/straysh/go_mnemonic/wordlist"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func Test_words(t *testing.T) {
	dict,err := LoadWordDict("japanese")
	if err!=nil {
		t.Error(err)
		t.FailNow()
	}
	assert.Equal(t, DictLength, dict.Len())
	fmt.Println(dict.PickIndex(101))
}
