package filex

import (
	"bytes"
	"path/filepath"
	"strings"
)

// binarySniffBytes is the number of leading bytes IsBinary inspects.
const binarySniffBytes = 8000

// Within reports whether path is root itself or lies inside root. The check
// is purely lexical: both paths are cleaned, ".." elements are resolved, and
// symbolic links are not followed, so a symlink inside root that points
// elsewhere still counts as within. Resolve both paths with
// filepath.EvalSymlinks first when links must be followed. An absolute path is
// never within a relative root, or vice versa, and comparison is
// case-sensitive on every platform.
func Within(root, path string) bool {
	rel, err := filepath.Rel(filepath.Clean(root), filepath.Clean(path))
	if err != nil || filepath.IsAbs(rel) {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// IsBinary reports whether data looks like binary content, that is whether
// its first 8000 bytes contain a NUL byte. Empty data is not binary.
func IsBinary(data []byte) bool {
	return bytes.IndexByte(data[:min(len(data), binarySniffBytes)], 0) >= 0
}
