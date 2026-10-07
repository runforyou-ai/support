package random

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestBytes(t *testing.T) {
	tests := []struct {
		n, want int
	}{
		{-5, 0},
		{0, 0},
		{1, 1},
		{32, 32},
		{1000, 1000},
	}
	for _, tt := range tests {
		got := Bytes(tt.n)
		if got == nil || len(got) != tt.want {
			t.Errorf("Bytes(%d) = %v (len %d), want non-nil len %d", tt.n, got, len(got), tt.want)
		}
	}
	seen := map[string]bool{}
	for range 1000 {
		s := string(Bytes(16))
		if seen[s] {
			t.Fatal("Bytes(16) repeated a value")
		}
		seen[s] = true
	}
}

func TestHex(t *testing.T) {
	tests := []struct {
		n, wantLen int
	}{
		{-1, 0},
		{0, 0},
		{1, 2},
		{16, 32},
	}
	for _, tt := range tests {
		got := Hex(tt.n)
		if len(got) != tt.wantLen {
			t.Errorf("len(Hex(%d)) = %d, want %d", tt.n, len(got), tt.wantLen)
		}
		if b, err := hex.DecodeString(got); err != nil || len(b) != max(tt.n, 0) {
			t.Errorf("Hex(%d) = %q does not decode to %d bytes: %v", tt.n, got, tt.n, err)
		}
	}
}

func TestBase64URL(t *testing.T) {
	tests := []struct {
		n, wantLen int
	}{
		{-1, 0},
		{0, 0},
		{1, 2},
		{16, 22},
		{32, 43},
	}
	for _, tt := range tests {
		got := Base64URL(tt.n)
		if len(got) != tt.wantLen {
			t.Errorf("len(Base64URL(%d)) = %d, want %d", tt.n, len(got), tt.wantLen)
		}
		b, err := base64.RawURLEncoding.DecodeString(got)
		if err != nil || len(b) != max(tt.n, 0) {
			t.Errorf("Base64URL(%d) = %q does not decode to %d bytes: %v", tt.n, got, tt.n, err)
		}
	}
}

func TestSHA256Hex(t *testing.T) {
	tests := []struct {
		in   []byte
		want string
	}{
		{nil, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{[]byte{}, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{[]byte("abc"), "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
	}
	for _, tt := range tests {
		if got := SHA256Hex(tt.in); got != tt.want {
			t.Errorf("SHA256Hex(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestToken(t *testing.T) {
	tests := []struct {
		n, wantLen int
	}{
		{-1, 0},
		{0, 0},
		{16, 22},
		{32, 43},
	}
	for _, tt := range tests {
		token, hash := Token(tt.n)
		if len(token) != tt.wantLen {
			t.Errorf("Token(%d) token len = %d, want %d", tt.n, len(token), tt.wantLen)
		}
		switch {
		case tt.wantLen == 0 && hash != "":
			t.Errorf("Token(%d) hash = %q, want empty", tt.n, hash)
		case tt.wantLen > 0 && hash != HashToken(token):
			t.Errorf("Token(%d) hash = %q, want HashToken(token) %q", tt.n, hash, HashToken(token))
		}
	}
}

func TestHashToken(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"abc", "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
	}
	for _, tt := range tests {
		if got := HashToken(tt.in); got != tt.want {
			t.Errorf("HashToken(%q) = %q, want %q", tt.in, got, tt.want)
		}
		if got := SHA256Hex([]byte(tt.in)); got != HashToken(tt.in) {
			t.Errorf("SHA256Hex and HashToken disagree for %q", tt.in)
		}
	}
}

func TestVerifyToken(t *testing.T) {
	token, hash := Token(32)
	tests := []struct {
		name, token, hash string
		want              bool
	}{
		{"match", token, hash, true},
		{"other token", token + "x", hash, false},
		{"other hash", token, HashToken("other"), false},
		{"empty hash", token, "", false},
	}
	for _, tt := range tests {
		if got := VerifyToken(tt.token, tt.hash); got != tt.want {
			t.Errorf("%s: VerifyToken = %v, want %v", tt.name, got, tt.want)
		}
	}
}
