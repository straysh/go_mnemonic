package tests

import (
	"testing"

	"fmt"
	"github.com/straysh/go_mnemonic"
)

func Test_create_random(t *testing.T) {
	//words := "advice owner gadget brick degree vanish coconut end among erupt gain once"
	//words := "advice owner gadget brick degree vanish coconut end among erupt gain oncd"
	m, err := mnemonic.NewMnemonic(mnemonic.ChineseSimplified)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = m.CreateRandom(128, "")
	if err != nil {
		fmt.Println(err)
	} else {
		t.Log( m.Mnemonic() )
	}
}

func Test_from_menmonic(t *testing.T) {
	var words = "advice owner gadget brick degree vanish coconut end among erupt gain once"
	var passpharse = ""
	m, err := mnemonic.NewMnemonic(mnemonic.Japanese)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = m.FromMnemonic(words, passpharse)
	if err != nil {
		t.Error(err)
	} else {
		t.Log( m.Mnemonic() )
	}

}
