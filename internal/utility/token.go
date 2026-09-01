package utility

import "math/rand/v2"

const (
	TokenAlphabetLatinLowercase            = "abcdefghijklmnopqrstuvwxyz"
	TokenAlphabetLatinUppercase            = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	TokenAlphabetDigits                    = "0123456789"
	TokenAlphabetSelectSpecial             = "_-"
	TokenAlphabetLatinWithDigitsAndSpecial = TokenAlphabetLatinLowercase + TokenAlphabetLatinUppercase + TokenAlphabetDigits + TokenAlphabetSelectSpecial
)

func GenerateRandomToken(length int, alphabet string) string {
	token := make([]byte, length)

	for i := range length {
		token[i] = alphabet[rand.IntN(len(alphabet))]
	}

	return string(token)
}
