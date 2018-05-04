package mnemonic

import (
	"crypto/rand"
	"errors"
	"strings"
	"fmt"
	"crypto/sha512"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha256"
	"strconv"
	"github.com/straysh/go_mnemonic/utils"
)

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

func isValidEntropyBitsLength(entLength int) bool {
	return entLength % 32 == 0 && entLength >= 128 && entLength <= 256
}

func (m *Mnemonic) generateRandomEntropy(ent uint) ([]byte, error) {
	if !m.isEntropyValid(ent) {
		return nil, errors.New("values must be ENT > 128 and ENT % 32 == 0")
	}

	entropy := make([]byte, ent/8)
	_, err := rand.Read( entropy )
	if err != nil {
		return nil, err
	}

	return entropy, nil
}

func entropy2Checksum(ent []byte) string {
	hash := sha256.Sum256(ent)
	entropyBitsLength   := len(ent) * bitsInByte
	checkhsumBitsLength := entropyBitsLength / multiple

	bin := make([]byte, 0)
	for i:=0;i<len(hash);i++ {
		bufi := utils.IntTo8BitsArray( int(hash[i]) )
		bin = append(bin, bufi[:]...)
	}

	checksum := bin[:checkhsumBitsLength]
	return string(checksum)
}

func (m *Mnemonic) isMnemonicValid(mnemonic string) bool {
	words := strings.Fields(mnemonic)
	mnemonicBitsLength := len(words) * wordBits
	checksumBitsLength := mnemonicBitsLength % multiple
	entropyBitsLength := mnemonicBitsLength - checksumBitsLength

	if isValidEntropyBitsLength(entropyBitsLength) == false {
		return false
	}

	binWithChecksum := make([]byte, 0)
	for _,word := range words {
		wordIndex, ok := m.WordList.SeekWord(word)
		if !ok {
			return false
		}
		buf := utils.IntTo11BitsArray(wordIndex)
		binWithChecksum = append(binWithChecksum, buf[:]...)
	}
	if len(binWithChecksum) != mnemonicBitsLength {
		return false
	}

	bin := binWithChecksum[:entropyBitsLength]
	if len(bin) != entropyBitsLength {
		return false
	}

	return true
}

func (m *Mnemonic) pickMnemonic(entropy []byte) (string, error) {
	entropyBitsLength := len(entropy)*bitsInByte
	if isValidEntropyBitsLength(entropyBitsLength) == false {
		return "", fmt.Errorf("values must be 128<= ENT <= 256 and ENT %% 32 == 0, current is %v", entropyBitsLength)
	}

	checksumBitsLength := entropyBitsLength / 32
	pharseLength := (entropyBitsLength + checksumBitsLength) / 11

	bin := make([]byte, 0)
	for _, b:= range entropy {
		buf := utils.IntTo8BitsArray(int(b))
		bin = append(bin, buf[:]...)
	}
	checksum := entropy2Checksum(entropy)
	bin = append(bin, checksum...)
	if len(bin)%11 !=0 {
		return "", fmt.Errorf("invalid entropy checksum length %v %% 11 !==0", len(bin))
	}
	pharse := make([]string, pharseLength)
	var byteAsBinaryString string
	for i:=0;i<pharseLength;i++{
		startIndex := i * wordBits
		endIndex   := startIndex + wordBits
		byteAsBinaryString = string(bin[startIndex:endIndex])
		asInt64, err := strconv.ParseInt(byteAsBinaryString, 2, 64)
		if err != nil {
			return "", err
		}
		pharse[i] = m.WordList.PickIndex(asInt64)
	}
	// \x20:ASCII标准空格/\xa0:nbsp(non-breaking space)不间断空白符/\u3000:全角空白符
	// 日语 Join(pharse, "\u3000")
	if m.Language == Japanese {
		return strings.Join(pharse, "\u3000"), nil
	} else {
		return strings.Join(pharse, " "), nil
	}
}

func (m *Mnemonic) mnemonic2Entropy(mnemonic string) ([]byte, error){
	words := strings.Fields(mnemonic)
	mnemonicBitsLength := len(words) * wordBits
	checksumBitsLength := mnemonicBitsLength % multiple
	entropyBitsLength := mnemonicBitsLength - checksumBitsLength

	if isValidEntropyBitsLength(entropyBitsLength) == false {
		return nil, fmt.Errorf("values must be 128<= ENT <= 256 and ENT %% 32 == 0, current is %v", entropyBitsLength)
	}

	binWithChecksum := make([]byte, 0)
	for _,word := range words {
		wordIndex, ok := m.WordList.SeekWord(word)
		if !ok {
			return nil, fmt.Errorf("invalid word {%v} in wordlist", word)
		}
		buf := utils.IntTo11BitsArray(wordIndex)
		binWithChecksum = append(binWithChecksum, buf[:]...)
	}
	if len(binWithChecksum) != mnemonicBitsLength {
		return nil, fmt.Errorf("mnemonicBitsLength should be %v, current %v, mnemonic: [%v]", mnemonicBitsLength, len(binWithChecksum), mnemonic)
	}

	bin := binWithChecksum[:entropyBitsLength]
	if len(bin) != entropyBitsLength {
		return nil, fmt.Errorf("entropyBitsLength should be %v, current %v, mnemonic: [%v]", entropyBitsLength, len(bin), mnemonic)
	}

	ent := make([]byte, entropyBitsLength / bitsInByte)
	var byteAsBinaryString string
	for i:=0;i<len(ent);i++ {
		startIndex := i * bitsInByte
		endIndex   := startIndex + bitsInByte
		if endIndex >= len(bin)-1 {
			byteAsBinaryString = string(bin[startIndex:])
		} else {
			byteAsBinaryString = string(bin[startIndex:endIndex])
		}
		asInt64, err := strconv.ParseInt(byteAsBinaryString, 2, 64)
		if err != nil {
			return nil, err
		}
		ent[i] = byte(asInt64)
	}

	checksum := binWithChecksum[entropyBitsLength:]
	if len(checksum) != checksumBitsLength{
		return nil, fmt.Errorf("checksum expected %d character, but got %d", entropyBitsLength, len(checksum))
	}
	if string(checksum)!=entropy2Checksum(ent) {
		return nil, fmt.Errorf("invalid mnemonic: invalid checksum")
	}
	return ent, nil
}

//NewSeed Based on a code (word list) returns the seed (hex bytes)
func NewSeed(mnecmonic, passphrase string) []byte {
	return pbkdf2.Key([]byte(mnecmonic), []byte("mnemonic"+passphrase), 2048, 64, sha512.New)
}