package str

import (
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

// IsHTTPURL reports whether s is an absolute http or https URL, such as a page
// address opened in a browser. Unlike IsHTTPOrigin, the scheme and host are
// matched case-insensitively. The authority after "://" must be a host with an
// optional port following the rules of IsHTTPOrigin, so user information,
// percent-encoded hosts and empty or out-of-range ports are rejected; a path,
// query and fragment are allowed. Surrounding whitespace is not trimmed.
func IsHTTPURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil || (!strings.EqualFold(u.Scheme, "http") && !strings.EqualFold(u.Scheme, "https")) {
		return false
	}
	rest, ok := strings.CutPrefix(s[len(u.Scheme):], "://")
	if !ok {
		return false
	}
	authority := rest
	if i := strings.IndexAny(rest, "/?#"); i >= 0 {
		authority = rest[:i]
	}
	authority, ok = lowerASCII(authority)
	return ok && isHostPort(authority)
}

// IsHTTPOrigin reports whether s is a web origin in serialized form: a
// lowercase "http://" or "https://" scheme followed by a lowercase host and an
// optional port, with nothing after them. The host is an ASCII DNS name of
// letters, digits and hyphens whose dot-separated labels are 1 to 63
// characters long and do not start or end with a hyphen, without a trailing
// dot; a dotted-decimal IPv4 address, required when the last label is
// numeric; or an IPv6 address in brackets without a zone. Internationalized
// names must be given in their "xn--" form. The port, when present, is a
// decimal number from 1 to 65535 without leading zeros; default ports are
// allowed. A trailing slash, path, query, fragment, user information and any
// other character are rejected.
func IsHTTPOrigin(s string) bool {
	rest, ok := strings.CutPrefix(s, "https://")
	if !ok {
		rest, ok = strings.CutPrefix(s, "http://")
	}
	return ok && isHostPort(rest)
}

// IsHTTPBaseURL reports whether s is an absolute http or https URL suitable as
// an API base address that request paths are appended to. It follows
// IsHTTPURL and additionally rejects '?' and '#' anywhere in s.
func IsHTTPBaseURL(s string) bool {
	return IsHTTPURL(s) && !strings.ContainsAny(s, "?#")
}

// lowerASCII returns s with ASCII letters lowercased, or false when s contains
// a non-ASCII byte.
func lowerASCII(s string) (string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return "", false
		}
	}
	return strings.ToLower(s), true
}

// isHostPort reports whether s is a lowercase host with an optional port, as
// described for IsHTTPOrigin.
func isHostPort(s string) bool {
	host, port := s, ""
	if strings.HasPrefix(s, "[") {
		end := strings.IndexByte(s, ']')
		if end < 0 {
			return false
		}
		host, port = s[1:end], s[end+1:]
		addr, err := netip.ParseAddr(host)
		if err != nil || !addr.Is6() || addr.Zone() != "" || strings.ToLower(host) != host {
			return false
		}
	} else {
		if i := strings.LastIndexByte(s, ':'); i >= 0 {
			host, port = s[:i], s[i:]
		}
		if !isHostName(host) {
			return false
		}
	}
	if port == "" {
		return true
	}
	digits, ok := strings.CutPrefix(port, ":")
	if !ok || digits == "" || digits[0] == '0' || len(digits) > 5 || strings.Trim(digits, "0123456789") != "" {
		return false
	}
	n, err := strconv.Atoi(digits)
	return err == nil && n <= 65535
}

// isHostName reports whether s is a lowercase ASCII DNS name, or a
// dotted-decimal IPv4 address when its last label is numeric.
func isHostName(s string) bool {
	if s == "" || len(s) > 253 {
		return false
	}
	labels := strings.Split(s, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' ||
			strings.Trim(label, "abcdefghijklmnopqrstuvwxyz0123456789-") != "" {
			return false
		}
	}
	// A numeric last label makes the whole host an IPv4 address, as in browsers.
	if strings.Trim(labels[len(labels)-1], "0123456789") == "" {
		addr, err := netip.ParseAddr(s)
		return err == nil && addr.Is4()
	}
	return true
}

// emailLocalChars are the characters allowed between the dots of an email
// local part: the RFC 5322 atext set.
const emailLocalChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+/=?^_`{|}~-"

// IsEmail reports whether s is a plain email address that common mail servers
// deliver, such as "taylor@example.com". The local part is at most 64
// characters: runs of ASCII letters, digits and the symbols
// !#$%&'*+/=?^_`{|}~- joined by single dots. The domain is a DNS name of at
// least two labels, matched case-insensitively and following the host name
// rules of IsHTTPOrigin, whose last label is not numeric. The whole address is
// at most 254 characters. IP address literals, single-label domains, display
// names, angle brackets, quoted local parts, non-ASCII addresses (which need
// SMTPUTF8 support to deliver) and surrounding whitespace are rejected. The
// domain is not checked against DNS.
func IsEmail(s string) bool {
	local, domain, found := strings.Cut(s, "@")
	if !found || len(s) > 254 || local == "" || len(local) > 64 {
		return false
	}
	for part := range strings.SplitSeq(local, ".") {
		if part == "" || strings.Trim(part, emailLocalChars) != "" {
			return false
		}
	}
	domain, ok := lowerASCII(domain)
	if !ok || !strings.Contains(domain, ".") {
		return false
	}
	last := domain[strings.LastIndexByte(domain, '.')+1:]
	return strings.Trim(last, "0123456789") != "" && isHostName(domain)
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
