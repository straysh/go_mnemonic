// [Specs] https://github.com/bgadrian/go-mnemonic/tree/master/bip39
// [Article] https://github.com/tyler-smith/go-bip39/blob/master/bip39.go
// [Article] https://github.com/upamune/bip39

// CS = ENT / 32
// MS = (ENT + CS) / 11
// the initial entropy length (ENT)
// mnemonic sentence (MS)
// checksum bits length (CS)
//|  ENT  | CS | ENT+CS |  MS  |
//+-------+----+--------+------+
//|  128  |  4 |   132  |  12  |
//|  160  |  5 |   165  |  15  |
//|  192  |  6 |   198  |  18  |
//|  224  |  7 |   231  |  21  |
//|  256  |  8 |   264  |  24  |
package mnemonic

import (
	"fmt"
	"github.com/straysh/go_mnemonic/wordlist"
)

type DictLang int
func (dict DictLang) String() string {
	switch dict {
	case English: return "english"
	case ChineseSimplified: return "chinese_simplified"
	case ChineseTraditional: return "chinese_traditional"
	case French: return "french"
	case Italian: return "italian"
	case Japanese: return "japanese"
	case Spanish: return "spanish"
	default:
		return "unknown word language"
	}
}
const (
	English DictLang = iota
	ChineseSimplified
	ChineseTraditional
	French
	Italian
	Japanese
	Spanish
)
type WordListInterface interface{
	PickIndex(index int64) string
	SeekWord(word string) (int, bool)
}
type Mnemonic struct {
	Language DictLang
	WordList WordListInterface
	entropy []byte
	passpharse string
	mnemonic string
}


const bitsInByte = 8
const wordBits = 11
const multiple = 32

func NewMnemonic(lang DictLang) (*Mnemonic, error) {
	if lang<English || lang>Spanish {
		return nil, fmt.Errorf("NewMnemonic(lang DictLang): %d, %s", lang, lang)
	}
	m := &Mnemonic{Language: lang}
	m.WordList,_ = wordlist.LoadWordDict(lang.String()) // DictLang 已经校验过,此处的错误可以忽略

	return m, nil
}

func (m *Mnemonic) CreateRandom(wordLen int, passpharse string) (*Mnemonic, error) {
	if wordLen!=12 && wordLen!=15 && wordLen!=18 && wordLen!=21 && wordLen!=24 {
		return nil, fmt.Errorf("invalid wordLen: %d, must be 12 or 15 or 18 or 21 or 24", wordLen)
	}
	mnemonicBitsLength := wordLen * wordBits
	checksumBitsLength := mnemonicBitsLength % multiple
	entSize := mnemonicBitsLength - checksumBitsLength
	ent, err := m.generateRandomEntropy( uint(entSize) )
	if err != nil {
		return nil, err
	}

	return m.FromEntropy(ent, passpharse)
}


func (m *Mnemonic) FromEntropy(ent []byte, passpharse string) (*Mnemonic, error) {
	mnemonic, err := m.pickMnemonic(ent)
	if err != nil {
		return nil, err
	}

	m.entropy = ent
	m.passpharse = passpharse
	m.mnemonic = mnemonic

	return m, nil
}

func (m *Mnemonic) FromMnemonic(mnemonic string, passpharse string) (*Mnemonic, error) {
	entropy, err := m.mnemonic2Entropy(mnemonic)
	if err != nil {
		return nil, err
	}
	m.entropy = entropy
	m.passpharse = passpharse
	m.mnemonic = mnemonic

	return m, nil
}

func (m *Mnemonic) isEntropyValid(ent uint) bool {
	return ent % 32 == 0 && ent >= 128 && ent <= 256
}

func (m *Mnemonic) IsValid(mnemonic string) bool {
	return m.isMnemonicValid(mnemonic)
}

func (m *Mnemonic) Mnemonic() string {
	return m.mnemonic
}

//func (m *Mnemonic) GetEntropy() []byte {
//	return m.entropy
//}

func (m *Mnemonic) Seed() []byte {
	seed := NewSeed(m.mnemonic, m.passpharse)

	return seed
}

func (m *Mnemonic) Passpharse() string {
	return m.passpharse
}