package validate

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

// zoneDirs returns the zone directories time.LoadLocation searches, in order:
// $ZONEINFO, read once like the time package does, then the system
// directories. Zip archives and embedded zone data are searched after them.
var zoneDirs = sync.OnceValue(func() []string {
	dirs := []string{"/usr/share/zoneinfo", "/usr/share/lib/zoneinfo", "/usr/lib/locale/TZ", "/etc/zoneinfo"}
	if env := os.Getenv("ZONEINFO"); env != "" {
		dirs = append([]string{env}, dirs...)
	}
	return dirs
})

// Timezone reports whether name is an IANA time zone name, such as
// "Asia/Shanghai" or "UTC", that time.LoadLocation can load. Names are
// case-sensitive on every platform: when the zone is read from a directory,
// each element of name must match the stored file and directory names
// exactly, so "asia/shanghai" is rejected even where a case-insensitive file
// system lets time.LoadLocation open it. It returns false for "" and "Local",
// which time.LoadLocation maps to UTC and the system zone.
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
	// The first directory whose file parses as zone data is the one
	// time.LoadLocation used; zip archives and embedded data are
	// case-sensitive and need no check.
	for _, dir := range zoneDirs() {
		data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(name)))
		if err != nil {
			continue
		}
		if _, err := time.LoadLocationFromTZData(name, data); err != nil {
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
