package utils

import (
	"crypto/sha1"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/pbkdf2"
)

func TestRandomCharString(t *testing.T) {
	funcAssert := assert.New(t)
	// Test-1 : RandomCharString() should return string with expected length
	length := 32
	token := RandomCharString(length)
	funcAssert.Equal(length, len(token))
	// Test-2 : RandomCharString() should output random string, new generated string should not be equal to the previous one
	newToken := RandomCharString(length)
	funcAssert.NotEqual(token, newToken)
}

func TestComputeHmac256(t *testing.T) {
	funcAssert := assert.New(t)
	// Test-1 : ComputeHmac256() should return the right Hmac256 string based on a derived key
	var derivedKey = pbkdf2.Key([]byte("secret"), []byte("salt"), 4096, 32, sha1.New)
	var message = "hello world"
	var expectedHmac = "5r32q7W+0hcBnqzQwJJUDzVGoVivXGSodTcHSqG/9Q8="
	hmac := ComputeHmac256(message, derivedKey)
	funcAssert.Equal(hmac, expectedHmac)
}
