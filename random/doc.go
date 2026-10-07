// Package random generates cryptographically secure random values and the
// hashes commonly stored alongside them.
//
// All functions read from crypto/rand, which never fails on supported
// platforms, so they return no error. A size of zero or less yields empty
// results.
//
// Example:
//
//	key := random.Hex(16)               // 32 hex characters
//	nonce := random.Base64URL(16)       // 22 URL-safe characters, no padding
//	token, hash := random.Token(32)     // hand out token, store hash
//	ok := random.HashToken(token) == hash
package random
