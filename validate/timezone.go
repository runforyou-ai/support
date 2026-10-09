package validate

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// zoneDirs are the system zone directories time.LoadLocation searches after
// $ZONEINFO, in order.
var zoneDirs = []string{"/usr/share/zoneinfo", "/usr/share/lib/zoneinfo", "/usr/lib/locale/TZ", "/etc/zoneinfo"}

// Timezone reports whether name is an IANA time zone name, such as
// "Asia/Shanghai" or "UTC", that time.LoadLocation can load. Names are
// case-sensitive on every platform: "asia/shanghai" is rejected even where a
// case-insensitive file system would let time.LoadLocation open it. It returns
// false for "" and "Local", which time.LoadLocation maps to UTC and the system
// zone.
//
// Zone data is read from the system at run time. Programs that must validate
// zones on systems without a zone database, such as minimal containers or
// Windows hosts without Go installed, should embed one by importing
// _ "time/tzdata" in their main package.
func Timezone(name string) bool {
	if name == "" || name == "Local" {
		return false
	}
	if _, err := time.LoadLocation(name); err != nil {
		return false
	}
	dirs := zoneDirs
	if env := os.Getenv("ZONEINFO"); env != "" {
		dirs = append([]string{env}, dirs...)
	}
	// The first directory holding the zone file decides, as in time.LoadLocation;
	// each element of name must match a directory entry exactly. Zip archives
	// and embedded zone data are case-sensitive and need no check.
	for _, dir := range dirs {
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(name))); err != nil {
			continue
		}
		for part := range strings.SplitSeq(name, "/") {
			entries, err := os.ReadDir(dir)
			if err != nil || !slices.ContainsFunc(entries, func(e os.DirEntry) bool { return e.Name() == part }) {
				return false
			}
			dir = filepath.Join(dir, part)
		}
		return true
	}
	return true
}
