package random

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// Bytes returns n cryptographically secure random bytes. It returns an empty
// slice when n is zero or negative.
func Bytes(n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	b := make([]byte, n)
	_, _ = rand.Read(b) // crypto/rand.Read never returns an error.
	return b
}

// Hex returns nBytes random bytes encoded as lowercase hexadecimal, so the
// result has 2*nBytes characters. It returns "" when nBytes is zero or
// negative.
func Hex(nBytes int) string {
	return hex.EncodeToString(Bytes(nBytes))
}

// Base64URL returns nBytes random bytes encoded with unpadded URL-safe base64
// (base64.RawURLEncoding), suitable for tokens, secrets and nonces in URLs and
// headers. It returns "" when nBytes is zero or negative.
func Base64URL(nBytes int) string {
	return base64.RawURLEncoding.EncodeToString(Bytes(nBytes))
}

// SHA256Hex returns the lowercase hexadecimal SHA-256 digest of data.
func SHA256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// Token returns a random token of nBytes bytes encoded as unpadded URL-safe
// base64, together with HashToken(token). Hand the token to the client and
// store only the hash; verify a presented token by comparing its HashToken
// with the stored hash. Token(32) yields a 256-bit token. Both results are ""
// when nBytes is zero or negative.
func Token(nBytes int) (token, hash string) {
	if nBytes <= 0 {
		return "", ""
	}
	token = Base64URL(nBytes)
	return token, HashToken(token)
}

// HashToken returns the lowercase hexadecimal SHA-256 digest of the token
// string itself (not of the decoded random bytes); it equals
// SHA256Hex([]byte(token)).
func HashToken(token string) string {
	return SHA256Hex([]byte(token))
}
