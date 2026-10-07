package data

import (
	"errors"
	"math/big"
	"net"
	"net/url"
	"testing"
	"time"
)

type tier int

func TestQuery(t *testing.T) {
	n := 5
	var nilPtr *int
	inner := m{"b": 1}
	list := []any{"x"}
	empty := m{}
	pp := &n
	seven := 7
	var boxed any = &seven
	var boxedBytes any = []byte("hi")
	var boxedNil any = (*int)(nil)
	ip := net.IP{127, 0, 0, 1}
	tests := []struct {
		name string
		in   m
		want string
	}{
		{"nil", nil, ""},
		{"flat sorted", m{"b": "2", "a": 1}, "a=1&b=2"},
		{"nested", m{"a": m{"b": "c"}}, "a%5Bb%5D=c"},
		{"list", m{"list": []any{"x", "y"}}, "list%5B0%5D=x&list%5B1%5D=y"},
		{"deep", m{"a": []any{m{"b": 1}}}, "a%5B0%5D%5Bb%5D=1"},
		{"escape", m{"q": "a b&c=d", "k y": "é"}, "k+y=%C3%A9&q=a+b%26c%3Dd"},
		{"bools", m{"t": true, "f": false}, "f=0&t=1"},
		{"skip nil", m{"a": nil, "b": nilPtr, "c": 1}, "c=1"},
		{"empty containers", m{"a": m{}, "b": []any{}, "c": 1}, "c=1"},
		{"bytes", m{"b": []byte("hi")}, "b=hi"},
		{"pointer", m{"p": &n}, "p=5"},
		{"pointer to pointer", m{"p": &pp}, "p=5"},
		{"pointer to map", m{"a": &inner}, "a%5Bb%5D=1"},
		{"pointer to slice", m{"a": &list}, "a%5B0%5D=x"},
		{"pointer to empty map", m{"a": &empty, "c": 1}, "c=1"},
		{"pointer stringer", m{"n": big.NewInt(42)}, "n=42"},
		{"pointer through interface", m{"v": &boxed}, "v=7"},
		{"bytes through interface", m{"b": &boxedBytes}, "b=hi"},
		{"nil pointer through interface", m{"p": &boxedNil, "c": 1}, "c=1"},
		{"stringer container pointer", m{"ip": &ip}, "ip=127.0.0.1"},
		{"stringer container", m{"ip": ip}, "ip=127.0.0.1"},
		{"pointer error", m{"e": &url.Error{Op: "Get", URL: "u", Err: errors.New("x")}}, "e=Get+%22u%22%3A+x"},
		{"numbers", m{"i": int8(-3), "u": uint(4), "f": 1.5, "g": float32(0.25)}, "f=1.5&g=0.25&i=-3&u=4"},
		{"named kinds", m{"t": tier(2), "l": level("hi"), "b": namedBool(true)}, "b=1&l=hi&t=2"},
		{"stringer", m{"d": time.Second}, "d=1s"},
		{"error", m{"e": errors.New("boom")}, "e=boom"},
		{"typed containers", m{"s": []int{1}, "m": map[string]string{"k": "v"}}, "m%5Bk%5D=v&s%5B0%5D=1"},
		{"fallback", m{"c": complex(1, 2)}, "c=%281%2B2i%29"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Query(tt.in); got != tt.want {
				t.Errorf("Query = %q; want %q", got, tt.want)
			}
		})
	}
}

type namedBool bool
