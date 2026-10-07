package str

import (
	"net/mail"
	"net/url"
	"strings"
	"unicode/utf8"
)

// IsHTTPURL reports whether s is an absolute http or https URL, such as a page
// address opened in a browser. The scheme is matched case-insensitively, the
// host must be non-empty and user information is rejected; a path, query and
// fragment are allowed. Surrounding whitespace is not trimmed.
func IsHTTPURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.IsAbs() && u.Host != "" && u.User == nil &&
		(strings.EqualFold(u.Scheme, "http") || strings.EqualFold(u.Scheme, "https"))
}

// IsHTTPOrigin reports whether s is a web origin: a lowercase "http://" or
// "https://" scheme followed by a non-empty host and an optional port, with
// nothing after them. It rejects a trailing slash, any path, query, fragment
// or user information, and any space; that is, the text after "://" contains
// none of '/', '?', '#', '@' or ' '.
func IsHTTPOrigin(s string) bool {
	u, err := url.Parse(s)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return false
	}
	return !strings.ContainsAny(strings.TrimPrefix(s, u.Scheme+"://"), "/?#@ ")
}

// IsHTTPBaseURL reports whether s is an absolute http or https URL suitable as
// an API base address. The host must be non-empty, user information and a
// non-empty query are rejected, and a path is allowed. s is parsed with
// url.ParseRequestURI, so '#' does not start a fragment: it is rejected in the
// host and kept as part of the path. Surrounding whitespace is not trimmed.
func IsHTTPBaseURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	return err == nil && u.IsAbs() && u.Host != "" && u.User == nil &&
		(u.Scheme == "http" || u.Scheme == "https") && u.RawQuery == "" && u.Fragment == ""
}

// IsEmail reports whether s is a bare RFC 5322 address as accepted by
// net/mail.ParseAddress, such as "taylor@example.com". Display names, angle
// brackets, comments, quoted local parts and surrounding whitespace are
// rejected. The domain need not contain a dot and is not checked against DNS.
func IsEmail(s string) bool {
	addr, err := mail.ParseAddress(s)
	return err == nil && strings.EqualFold(addr.Address, s)
}

// MaskEmail hides the local part of an email address except for its first
// character, replacing the rest with a fixed "***" so that its length is not
// revealed: "taylor@example.com" becomes "t***@example.com". The address is
// split at the first '@' and the domain is kept as is. It returns email
// unchanged when it contains no '@' or the local part is empty.
func MaskEmail(email string) string {
	local, domain, found := strings.Cut(email, "@")
	if !found || local == "" {
		return email
	}
	first, _ := utf8.DecodeRuneInString(local)
	return string(first) + "***@" + domain
}
