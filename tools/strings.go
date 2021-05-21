package tools

import (
	"encoding/base32"
	"math/rand"
	"strings"
)

/* Shuffle string */
func Shuffle(input string) (output string) {
	inRune := []rune(input)
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

/* XOR encrypt/decrypt */
func XOR(input, key string) (output string) {
	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ key[i%len(key)])
	}
	return output
}

/* Light string encoder */
func StringEncode(st, code string) string {
	return strings.TrimRight(base32.StdEncoding.EncodeToString([]byte(XOR(st, code))), "=")
}

/* Light string decoder */
func StringDecode(st, code string) (string, error) {
	str := st
	for (len(str) % 8) != 0 {
		str = str + "="
	}
	a, err := base32.StdEncoding.DecodeString(str)
	var b string = ""
	if err == nil {
		b = XOR(string(a), code)
	}
	return b, err
}
