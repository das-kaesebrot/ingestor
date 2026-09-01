package utility

import "math/rand/v2"

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

func GenerateRandomToken(length int, alphabet string) string {
	token := make([]byte, length)

	for i := range length {
		token[i] = alphabet[rand.IntN(len(alphabet))]
	}

	return string(token)
}
