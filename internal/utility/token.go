package utility

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	TokenAlphabetLatinLowercase            = "abcdefghijklmnopqrstuvwxyz"
	TokenAlphabetLatinUppercase            = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	TokenAlphabetDigits                    = "0123456789"
	TokenAlphabetSelectSpecial             = "_-"
	TokenAlphabetLatinWithDigitsAndSpecial = TokenAlphabetLatinLowercase + TokenAlphabetLatinUppercase + TokenAlphabetDigits + TokenAlphabetSelectSpecial
	// Characters which are ambiguous are excluded, such as I, l, and 1 and so on.
	// inspired from Nextcloud share ID generator: https://github.com/nextcloud/server/blob/master/lib/public/Security/ISecureRandom.php
	TokenAlphabetHumanReadable = "abcdefgijkmnopqrstwxyzABCDEFGHJKLMNPQRSTWXYZ23456789"
)

// https://stackoverflow.com/a/6878625
const (
	maxUint = ^uint(0)
	minUint = 0
	maxInt  = int(maxUint >> 1)
	minInt  = -maxInt - 1
)

func GenerateRandomToken(length int, alphabet string) (string, error) {
	token := make([]byte, length)

	if len(alphabet) > maxInt {
		return "", fmt.Errorf("Alphabet can't be bigger than int max val! maxInt=%v, alphabetLength=%v", maxInt, len(alphabet))
	}

	for i := range length {
		a, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		token[i] = alphabet[int(a.Int64())]
	}

	return string(token), nil
}
