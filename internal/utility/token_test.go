package utility

import (
	"strings"
	"testing"
)

var tokenGeneratorTests = []struct {
	length   int
	alphabet string
}{
	{10, TokenAlphabetDigits},
	{10, TokenAlphabetLatinLowercase},
	{10, TokenAlphabetLatinUppercase},
	{10, TokenAlphabetSelectSpecial},
	{10, TokenAlphabetLatinWithDigitsAndSpecial},
	{128, TokenAlphabetLatinWithDigitsAndSpecial},
	{16, TokenAlphabetHumanReadable},
}

func TestTokenGenerator(t *testing.T) {
	for _, tt := range tokenGeneratorTests {
		token := GenerateRandomToken(tt.length, tt.alphabet)

		if len(token) != tt.length {
			t.Errorf("Token length is not expected value! expected=%v, result=%v", tt.length, len(token))
		}
		remaining := strings.Trim(token, tt.alphabet)

		if len(remaining) != 0 {
			t.Errorf("Token contains characters not in alphabet! found=%v, alphabet=%v", remaining, tt.alphabet)
		}
	}
}
