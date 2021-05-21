package tools

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func SHA1HashPassword(str string) (string, error) {
	sha := sha1.New()
	_, err := sha.Write([]byte(str))
	if err != nil {
		return "", err
	}
	hashedPw := sha.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashedPw), nil
}

func SHA1hash(b []byte) [20]byte {
	return sha1.Sum(b)
}

func SHA256hash(b []byte) [32]byte {
	return sha256.Sum256(b)
}

func SHA512hash(b []byte) [64]byte {
	return sha512.Sum512(b)
}

func MD5hash(b []byte) [16]byte {
	return md5.Sum(b)
}

func MD5HashPassword(str string) (string, error) {
	md5 := md5.New()
	_, err := md5.Write([]byte(str))
	if err != nil {
		return "", err
	}
	hashedPw := md5.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashedPw), nil
}

func NoHashPassword(str string) (string, error) {
	return str, nil
}

func BCRYPTHashPassword(str string, cost int) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(str), cost)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

func MD5sum(st string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(st)))
}

// Return sha256 from string
func Gensha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
